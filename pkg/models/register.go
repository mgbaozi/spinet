package models

import (
	"k8s.io/klog/v2"
)

var registeredTypes *RegisteredTypes

type RegisteredTypes struct {
	Triggers  map[string]Trigger
	Apps      map[string]App
	Operators map[string]Operator
	Handlers  []Handler
}

func init() {
	//initialize static instance on load
	registeredTypes = &RegisteredTypes{
		Triggers:  make(map[string]Trigger),
		Apps:      make(map[string]App),
		Operators: make(map[string]Operator),
	}
}

func GetRegisteredTypes() *RegisteredTypes {
	return registeredTypes
}

func RegisterTrigger(trigger Trigger) {
	name := trigger.TriggerName()
	klog.V(2).Infof("Register trigger: %s", name)
	registeredTypes.Triggers[name] = trigger
}

func RegisterApp(app App) {
	name := app.AppName()
	klog.V(2).Infof("Register app: %s", name)
	registeredTypes.Apps[name] = app
}

func RegisterOperator(operator Operator) {
	name := operator.Name()
	klog.V(2).Infof("Register operator: %s", name)
	registeredTypes.Operators[name] = operator
}

func RegisterOperators(operators ...Operator) {
	for _, operator := range operators {
		RegisterOperator(operator)
	}
}

func RegisterHandler(handler Handler) {
	klog.V(2).Infof("Register handler: %s", handler.Type())
	registeredTypes.Handlers = append(registeredTypes.Handlers, handler)
}

func NewTrigger(name string, options map[string]interface{}) Trigger {
	trigger := registeredTypes.Triggers[name]
	return trigger.New(options)
}

func NewApp(name string, options map[string]interface{}) App {
	app := registeredTypes.Apps[name]
	return app.New(options)
}

func NewOperator(name string) Operator {
	operator := registeredTypes.Operators[name]
	return operator
}

func GetHandlers() []Handler {
	return registeredTypes.Handlers
}
