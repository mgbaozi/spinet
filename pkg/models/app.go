package models

import "k8s.io/klog/v2"

type CustomApp struct {
	Task
	Modes   []AppMode
	Options map[string]interface{}
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	return &CustomApp{
		Task:    custom.Task,
		Modes:   custom.Modes,
		Options: options,
	}
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

func (custom *CustomApp) Execute(ctx *Context, mode AppMode, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("Execute app %s failed with error %v", custom.Name, err)
		}
		klog.V(2).Infof("App %s finished", custom.Name)
	}()
	var res bool
	if res, err = processSteps(&custom.Context, custom.Inputs, string(AppModeInput)); err != nil || !res {
		return
	}
	return
}
