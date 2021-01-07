package models

import (
	"k8s.io/klog/v2"
	"strings"
)

var registeredTypes *RegisteredTypes

type RegisteredTypes struct {
	Triggers         map[string]Trigger
	Apps             map[string]App
	Operators        map[string]Operator
	BuildInVariables map[string]BuildInVariable
	Handlers         []Handler
}

// func init() {
// 	//initialize static instance on load
// 	registeredTypes = &RegisteredTypes{
// 		Triggers:       make(map[string]Trigger),
// 		Apps:           make(map[string]App),
// 		Operators:      make(map[string]Operator),
// 		BuildInVariables: make(map[string]BuildInVariable),
// 	}
// }

func GetRegisteredTypes() *RegisteredTypes {
	if registeredTypes == nil {
		registeredTypes = &RegisteredTypes{
			Triggers:         make(map[string]Trigger),
			Apps:             make(map[string]App),
			Operators:        make(map[string]Operator),
			BuildInVariables: make(map[string]BuildInVariable),
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

func RegisterOperator(operator Operator) {
	name := operator.Name()
	klog.V(2).Infof("Register operator: %s", name)
	GetRegisteredTypes().Operators[name] = operator
}

func RegisterOperators(operators ...Operator) {
	for _, operator := range operators {
		RegisterOperator(operator)
	}
}

func RegisterBuildInVariable(variable BuildInVariable) {
	name := variable.Name()
	klog.V(2).Infof("Register build-in variable: %s", name)
	GetRegisteredTypes().BuildInVariables[name] = variable
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

func NewOperator(name string) Operator {
	operator := GetRegisteredTypes().Operators[name]
	return operator
}

func NewBuildInVariable(name string, value interface{}) BuildInVariable {
	name = strings.ToLower(name)
	if variable, ok := GetRegisteredTypes().BuildInVariables[name]; ok {
		return variable.New(value)
	}
	return nil
}

func GetBuildInVariables() map[string]BuildInVariable {
	return GetRegisteredTypes().BuildInVariables
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
