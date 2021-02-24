package models

import (
	"k8s.io/klog/v2"
)

var registeredTypes *RegisteredTypes

type RegisteredTypes struct {
	Triggers map[string]Trigger
	Apps     map[string]App
	Handlers []Handler
}

func GetRegisteredTypes() *RegisteredTypes {
	if registeredTypes == nil {
		registeredTypes = &RegisteredTypes{
			Triggers: make(map[string]Trigger),
			Apps:     make(map[string]App),
		}
	}
	return registeredTypes
}

func RegisterTrigger(trigger Trigger) {
	name := trigger.TriggerName()
	klog.V(2).Infof("Register trigger: %s", name)
	GetRegisteredTypes().Triggers[name] = trigger
}

func RegisterApp(app App) {
	name := app.AppName()
	klog.V(2).Infof("Register app: %s", name)
	GetRegisteredTypes().Apps[name] = app
}

func RegisterHandler(handler Handler) {
	klog.V(2).Infof("Register handler: %s", handler.Type())
	GetRegisteredTypes().Handlers = append(GetRegisteredTypes().Handlers, handler)
}

func NewTrigger(name string, options map[string]interface{}) Trigger {
	trigger := GetRegisteredTypes().Triggers[name]
	return trigger.New(options)
}

func NewApp(name string, options map[string]interface{}) (App, error) {
	app := GetRegisteredTypes().Apps[name]
	return app.New(options), nil
}

func GetHandlers() []Handler {
	return GetRegisteredTypes().Handlers
}

func GetApps() []App {
	var res []App
	for _, app := range GetRegisteredTypes().Apps {
		res = append(res, app)
	}
	return res
}
