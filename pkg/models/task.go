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
	AppData map[string]interface{}
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
	context := NewContext()
	if dictionary != nil {
		context.Dictionary = dictionary
	}
	return context
}

type Task struct {
	Meta
	Triggers   []Trigger
	Inputs     []Step
	Conditions []Condition
	Outputs    []Step
	// FIXME: context need refresh before each execution
	// TODO: versioned context
	Context          Context
	originDictionary map[string]interface{}
}

// set origin dictionary of task, it will set to context before every execution
func (task *Task) SetDictionary(dictionary map[string]interface{}) {
	if dictionary == nil {
		dictionary = make(map[string]interface{})
	}
	task.originDictionary = dictionary
}

func (task *Task) Start() {
	task.Context.Meta = task.Meta
	//FIXME: only the first trigger will active
	//FIXME: check if task has trigger
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

func (task *Task) prepare() {
	if task.originDictionary == nil {
		task.originDictionary = make(map[string]interface{})
	}
	task.Context = NewContextWithDictionary(task.originDictionary)
}

func ProcessMapper(ctx *Context, mapper Mapper, data interface{}) {
	for key, value := range mapper {
		if v, err := value.Extract(ctx.Dictionary, data); err == nil {
			ctx.Dictionary[key] = v
			klog.V(2).Infof("Set value %v to ctx.dictionary with key %s", v, key)
		}
	}
}

func (task *Task) processConditions() (bool, error) {
	return ProcessConditions(NewOperator("and"), task.Conditions, task.Context.Dictionary, nil)
}

func (task *Task) Execute() (res bool, err error) {
	defer func() {
		if err != nil {
			//FIXME: error description
			klog.V(3).Infof("Process conditions of task %s failed with error %v, skip output...", task.Name, err)
		} else if !res {
			klog.V(3).Infof("Conditions are not true, skip output...")
		}
		klog.V(2).Infof("Task %s finished", task.Name)
	}()
	klog.V(2).Infof("Running task %s", task.Name)
	task.prepare()
	if res, err = processSteps(&task.Context, task.Inputs, string(AppModeInput)); err != nil || !res {
		return
	}
	if res, err = task.processConditions(); err != nil || !res {
		return
	}
	return processSteps(&task.Context, task.Outputs, string(AppModeOutPut))
	// TODO: add output validator
	// TODO: add task validator
}
