package models

import (
	"github.com/mgbaozi/spinet/pkg/operators"
	"github.com/mgbaozi/spinet/pkg/values"
	"k8s.io/klog/v2"
	"reflect"
	"time"
)

const DefaultNamespace = "default"

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
	Conditions []values.Value
	Outputs    []Step
	Aggregator Mapper
	// TODO: versioned context
	Context          Context
	OriginDictionary map[string]values.Value
}

// set origin dictionary of task, it will set to context before every execution
func (task *Task) SetDictionary(dictionary map[string]interface{}) {
	task.OriginDictionary = make(map[string]values.Value)
	if dictionary == nil {
		dictionary = make(map[string]interface{})
	}
	for key, item := range dictionary {
		task.OriginDictionary[key] = values.Parse(item)
	}
}

func (task *Task) Start() {
	if len(task.Triggers) == 0 {
		klog.Errorf("Task %s.%s has no triggers, will never be called", task.Name, task.Namespace)
		return
	}
	task.init()
	//TODO: handle interrupt, timeout, heartbeats
	cases := make([]reflect.SelectCase, len(task.Triggers))
	for i, trigger := range task.Triggers {
		ch := trigger.Triggered(&task.Context)
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	for {
		_, recv, _ := reflect.Select(cases)
		hookData := recv.Interface().(map[string]interface{})
		klog.V(4).Infof("trigger with data %v, %s", hookData, recv.Kind())
		task.prepare(hookData)
		task.Execute()
		klog.V(9).Infof("%s", task.Context.trace.String())
	}
}

func (task *Task) init() {
	if task.OriginDictionary == nil {
		task.OriginDictionary = make(map[string]values.Value)
	}
	task.Context = NewContext()
	task.Context.Meta = task.Meta
}

func (task *Task) prepare(data map[string]interface{}) {
	klog.V(4).Infof("Prepare task with initial data: %v", data)
	dictionary := make(map[string]interface{})
	for key, item := range task.OriginDictionary {
		dictionary[key], _ = item.Extract(task.Context.BuildIn)
	}
	for key, item := range data {
		dictionary[key] = item
	}
	task.Context = NewContextWithDictionary(dictionary)
	task.Context.Meta = task.Meta
}

func (task *Task) processConditions() (bool, error) {
	//FIXME
	return ProcessConditions(task.Context, operators.New("and"), task.Conditions)
}

func (task *Task) Execute() (result map[string]interface{}, err error) {
	var res bool
	defer func() {
		if err != nil {
			//FIXME: error description
			klog.V(3).Infof("Task %s execute failed with error %v...", task.Name, err)
		} else if !res {
			klog.V(3).Infof("Conditions are not true, skip steps...")
		}
		task.Context.Trace(err == nil, "task finish", res)
		klog.V(2).Infof("Task %s finished", task.Name)
	}()
	task.Context.Trace(true, "task start", nil)
	klog.V(2).Infof("Running task %s", task.Name)
	if res, err = ProcessSteps(task.Context.Sub(string(TaskProgressInput), nil), task.Inputs, nil); err != nil || !res {
		return
	}
	if res, err = task.processConditions(); err != nil || !res {
		return
	}
	if res, err = ProcessSteps(task.Context.Sub(string(TaskProgressOutput), nil), task.Outputs, nil); err != nil || !res {
		return
	}
	// TODO: add output validator
	// TODO: add task validator
	// TODO: this return value is task data, save to history
	result = ProcessMapper(task.Aggregator, task.Context.Dictionary)
	return result, nil
}
