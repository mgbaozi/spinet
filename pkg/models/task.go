package models

import (
	"k8s.io/klog/v2"
	"time"
)

type Meta struct {
	Name      string
	Namespace string
}

type Status struct {
	Phase     string
	StartTime *time.Time
	EndTime   *time.Time
}

type Context struct {
	Meta       Meta
	Status     Status
	Dictionary map[string]interface{}
	// NamedAppData    map[string]interface{}
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
	Meta
	Triggers   []Trigger
	Inputs     []Input
	Conditions []Condition
	Outputs    []Output
	// FIXME: context need refresh before each execution
	// TODO: versioned context
	Context Context
}

func (task *Task) Start() {
	task.Context.Meta = task.Meta
	for _, trigger := range task.Triggers {
		for {
			select {
			case <-trigger.Triggered(&task.Context):
				task.Execute()
				//TODO: handle interrupt, timeout, heartbeats
			}
		}
	}

}

func ProcessMapper(ctx *Context, mapper Mapper, data interface{}) {
	for key, value := range mapper {
		if v, err := value.Extract(ctx.Dictionary, data); err == nil {
			ctx.Dictionary[key] = v
		}
	}
}

func (task *Task) processConditions() (bool, error) {
	return ProcessConditions(NewOperator("and"), task.Conditions, task.Context.Dictionary, nil)
}

func (task *Task) Execute() {
	defer func() {
		klog.V(2).Infof("Task %s finished", task.Name)
	}()
	klog.V(2).Infof("Running task %s", task.Name)
	res, err := processInputs(&task.Context, task.Inputs)
	if err != nil {
		klog.V(3).Infof("Condition execute failed with error: %v", err)
	}
	if !res {
		klog.V(3).Infof("Conditions in inputs are not true, skip output")
		return
	}
	res, err = task.processConditions()
	if err != nil {
		klog.V(3).Infof("Process conditions of task %s failed with error %v, skip output...", task.Name, err)
		return
	}
	if !res {
		klog.V(3).Infof("Conditions of task are not true, skip output...")
		return
	}
	for _, output := range task.Outputs {
		app := output.App
		var data interface{}
		_ = app.Execute(AppModeOutPut, &task.Context, &data)
		// TODO: add output validator
		// TODO: process output mapper
	}
	// TODO: add task validator
}
