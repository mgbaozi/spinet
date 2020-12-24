package models

import (
	"k8s.io/klog/v2"
	"reflect"
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

type Task struct {
	Meta
	Triggers   []Trigger
	Inputs     []Step
	Conditions []Condition
	Outputs    []Step
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
	if len(task.Triggers) == 0 {
		klog.Errorf("Task %s.%s has no triggers, will never be called", task.Name, task.Namespace)
		return
	}
	//TODO: handle interrupt, timeout, heartbeats
	cases := make([]reflect.SelectCase, len(task.Triggers))
	for i, trigger := range task.Triggers {
		ch := trigger.Triggered(task.Context)
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	for {
		_, _, _ = reflect.Select(cases)
		task.Execute()
		klog.V(9).Infof("%s", task.Context.Trace.String())
	}
}

func (task *Task) prepare() {
	if task.originDictionary == nil {
		task.originDictionary = make(map[string]interface{})
	}
	task.Context = NewContextWithDictionary(task.originDictionary)
}

func ProcessMapper(ctx Context, mapper Mapper, data interface{}) {
	for key, value := range mapper {
		//TODO: super data
		if v, err := value.Extract(ctx); err == nil {
			ctx.Dictionary[key] = v
			klog.V(2).Infof("Set value %v to ctx.dictionary with key %s", v, key)
		}
	}
}

func (task *Task) processConditions() (bool, error) {
	//FIXME
	return ProcessConditions(task.Context, NewOperator("and"), task.Conditions)
}

func (task *Task) Execute() (res bool, err error) {
	defer func() {
		if err != nil {
			//FIXME: error description
			klog.V(3).Infof("Process conditions of task %s failed with error %v, skip output...", task.Name, err)
		} else if !res {
			klog.V(3).Infof("Conditions are not true, skip output...")
		}
		task.Context.Trace.Push(err == nil, "task finish", res)
		klog.V(2).Infof("Task %s finished", task.Name)
	}()
	task.Context.Trace.Push(true, "task start", nil)
	klog.V(2).Infof("Running task %s", task.Name)
	task.prepare()
	magic := map[string]interface{}{
		"__mode__": string(AppModeInput),
	}
	if res, err = processSteps(task.Context.Sub(string(AppModeInput), magic), task.Inputs); err != nil || !res {
		return
	}
	if res, err = task.processConditions(); err != nil || !res {
		return
	}
	magic = map[string]interface{}{
		"__mode__": string(AppModeInput),
	}
	return processSteps(task.Context.Sub(string(AppModeOutPut), magic), task.Outputs)
	// TODO: add output validator
	// TODO: add task validator
}
