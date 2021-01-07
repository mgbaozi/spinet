package models

import (
	"k8s.io/klog/v2"
	"reflect"
)

type CustomApp struct {
	Task
	AppOptions map[string]Value
	options    map[string]interface{}
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	app := &CustomApp{
		Task:       custom.Task,
		AppOptions: make(map[string]Value),
	}
	if app.OriginDictionary == nil {
		app.OriginDictionary = make(map[string]Value)
	}
	for key, item := range options {
		app.AppOptions[key] = ParseValue(item)
	}
	app.options = options
	return app
}

func (custom *CustomApp) AppName() string {
	return custom.Name
}

func (custom *CustomApp) Register() {
	RegisterApp(custom)
}

func (custom *CustomApp) Options() map[string]interface{} {
	return custom.options
}

func (custom *CustomApp) prepare(ctx Context) (err error) {
	for key, value := range custom.AppOptions {
		//TODO: super data
		custom.OriginDictionary[key] = value
	}
	dictionary := make(map[string]interface{})
	for key, item := range custom.OriginDictionary {
		dictionary[key], _ = item.Extract(ctx.MergedData())
	}
	custom.Context = NewContextWithDictionary(dictionary)
	return nil
}

func (custom *CustomApp) Execute(ctx Context, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("Execute app %s failed with error %v", custom.Name, err)
		}
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val.Elem().Set(reflect.ValueOf(custom.Context.Dictionary))
		}
		klog.V(2).Infof("App %s finished", custom.Name)
	}()
	if err := custom.prepare(ctx); err != nil {
		return err
	}
	var res bool
	if res, err = processSteps(custom.Context.Sub(string(TaskProgressInput), nil), custom.Inputs); err != nil || !res {
		return
	}
	if res, err = custom.processConditions(); err != nil || !res {
		return
	}
	if res, err = processSteps(custom.Context.Sub(string(TaskProgressOutput), nil), custom.Outputs); err != nil || !res {
		return
	}
	return
}
