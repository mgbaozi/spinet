package models

import "fmt"

type Context struct {
	Dictionary map[string]interface{}
	Data       map[string]interface{}
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
	}
}

type Task struct {
	Name     string
	Triggers []Trigger
	Inputs   []Input
	Outputs  []Output
	Context  Context
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
	for _, input := range task.Inputs {
		app := input.App
		fmt.Println("Running app:", app.Name())
		var data map[string]interface{}
		err := app.Execute(AppModeInput, &task.Context, &data)
		if err != nil {
			fmt.Println(err)
		}
		task.Context.Data[app.Name()] = data
		ProcessAppConditions(input.Conditions, &task.Context)
	}
	for _, output := range task.Outputs {
		app := output.App
		app.Execute(AppModeOutPut, &task.Context, nil)
	}
}
