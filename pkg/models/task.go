package models

import "fmt"

type Context struct {
	Dictionary map[string]interface{}
	AppData    map[string]interface{}
}

func (ctx *Context) GetVariable(name string) interface{} {
	return ctx.Dictionary[name]
}

func (ctx *Context) SetVariable(name, value string) {
	ctx.Dictionary[name] = value
}

func NewContext() Context {
	return Context{
		Dictionary: make(map[string]interface{}),
		AppData:    make(map[string]interface{}),
	}
}

func NewContextWithDictionary(dictionary map[string]interface{}) Context {
	if dictionary == nil {
		dictionary = make(map[string]interface{})
	}
	return Context{
		Dictionary: dictionary,
		AppData:    make(map[string]interface{}),
	}
}

type Task struct {
	Name       string
	Triggers   []Trigger
	Inputs     []Input
	Conditions []Condition
	Outputs    []Output
	Context    Context
}

func (task *Task) Start() {
	for _, trigger := range task.Triggers {
		for {
			select {
			case <-trigger.Triggered():
				task.Execute()
			}
		}
	}

}

func (task *Task) Execute() {
	var inputResults []interface{}
	for _, input := range task.Inputs {
		app := input.App
		fmt.Println("Running app:", app.Name())
		var data map[string]interface{}
		err := app.Execute(AppModeInput, &task.Context, &data)
		if err != nil {
			fmt.Println(err)
		}
		task.Context.AppData[app.Name()] = data
		res, err := ProcessAppConditions(app.Name(), input.Conditions, &task.Context)
		if err != nil {
			fmt.Println(err)
		}
		inputResults = append(inputResults, res)
	}
	if res, err := (And{}).Do(inputResults); err != nil || !res {
		fmt.Println("Conditions in inputs are not true, skip output...")
		return
	}
	if res, err := ProcessCommonConditions(task.Conditions, &task.Context); err != nil || !res {
		fmt.Println("Conditions of task are not true, skip output...")
		return
	}
	for _, output := range task.Outputs {
		app := output.App
		_ = app.Execute(AppModeOutPut, &task.Context, nil)
	}
}
