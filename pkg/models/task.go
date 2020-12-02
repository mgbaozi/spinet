package models

import (
	"k8s.io/klog/v2"
)

type Context struct {
	Dictionary map[string]interface{}
	// AppData    map[string]interface{}
	AppData []interface{}
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

func NewContextWithDictionary(dictionary map[string]interface{}) Context {
	if dictionary == nil {
		dictionary = make(map[string]interface{})
	}
	return Context{
		Dictionary: dictionary,
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

func (task *Task) processMapper(input *Input, data interface{}) {
	for key, value := range input.Mapper {
		if v, err := value.Extract(task.Context.Dictionary, data); err == nil {
			task.Context.Dictionary[key] = v
		}
	}
}

func (task *Task) processInputConditions(input *Input, appdata interface{}) (bool, error) {
	return ProcessConditions(NewOperator("and"), input.Conditions, task.Context.Dictionary, appdata)
}

func (task *Task) processConditions() (bool, error) {
	return ProcessConditions(NewOperator("and"), task.Conditions, task.Context.Dictionary, nil)
}

func (task *Task) Execute() {
	var inputResults []interface{}
	for _, input := range task.Inputs {
		app := input.App
		klog.V(2).Infof("Running app: %s", app.Name())
		var data map[string]interface{}
		err := app.Execute(AppModeInput, &task.Context, &data)
		if err != nil {
			klog.Errorf("Execute app failed: %v", err)
		}
		task.Context.AppData = append(task.Context.AppData, data)
		task.processMapper(&input, data)
		res, err := task.processInputConditions(&input, data)
		if err != nil {
			klog.Errorf("Process conditions of app %s failed: %v", err)
		}
		inputResults = append(inputResults, res)
	}
	if res, err := NewOperator("and").Do(inputResults); err != nil || !res {
		klog.V(2).Infof("Conditions in inputs are not true, skip output...")
		return
	}
	if res, err := task.processConditions(); err != nil || !res {
		klog.V(2).Infof("Conditions of task are not true, skip output...")
		return
	}
	for _, output := range task.Outputs {
		app := output.App
		_ = app.Execute(AppModeOutPut, &task.Context, nil)
	}
}
