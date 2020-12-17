package models

import "k8s.io/klog/v2"

type CustomApp struct {
	Task
	Modes []AppMode
	// Options map[string]interface{}
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	app := &CustomApp{
		Task:  custom.Task,
		Modes: custom.Modes,
	}
	if app.originDictionary == nil {
		app.originDictionary = make(map[string]interface{})
	}
	//TODO: merge options need parse Value
	for key, value := range options {
		app.originDictionary[key] = value
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

func (custom *CustomApp) prepare() {
	custom.Context = NewContextWithDictionary(custom.originDictionary)
}

func (custom *CustomApp) Execute(ctx *Context, mode AppMode, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("Execute app %s failed with error %v", custom.Name, err)
		}
		klog.V(2).Infof("App %s finished", custom.Name)
	}()
	custom.prepare()
	var res bool
	if res, err = processSteps(&custom.Context, custom.Inputs, string(AppModeInput)); err != nil || !res {
		return
	}
	if res, err = processSteps(&custom.Context, custom.Outputs, string(AppModeOutPut)); err != nil || !res {
		return
	}
	return
}
