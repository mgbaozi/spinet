package models

import (
	"k8s.io/klog/v2"
	"reflect"
)

type CustomApp struct {
	Task
	Modes   []AppMode
	Options map[string]Value
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	app := &CustomApp{
		Task:    custom.Task,
		Modes:   custom.Modes,
		Options: make(map[string]Value),
	}
	if app.originDictionary == nil {
		app.originDictionary = make(map[string]interface{})
	}
	for key, item := range options {
		app.Options[key] = ParseValue(item)
	}
	return app
}

func (custom *CustomApp) AppName() string {
	return custom.Name
}

func (custom *CustomApp) Register() {
	RegisterApp(custom)
}

func (custom *CustomApp) AppModes() []AppMode {
	return custom.Modes
}

func (custom *CustomApp) prepare(ctx *Context) (err error) {
	custom.Context = NewContextWithDictionary(custom.originDictionary)
	for key, value := range custom.Options {
		if custom.originDictionary[key], err = value.Extract(ctx.Dictionary, nil); err != nil {
			return
		}
	}
	return nil
}

func (custom *CustomApp) Execute(ctx *Context, mode AppMode, data interface{}) (err error) {
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
	if res, err = processSteps(&custom.Context, custom.Inputs, string(AppModeInput)); err != nil || !res {
		return
	}
	if res, err = custom.processConditions(); err != nil || !res {
		return
	}
	if res, err = processSteps(&custom.Context, custom.Outputs, string(AppModeOutPut)); err != nil || !res {
		return
	}
	return
}
