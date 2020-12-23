package models

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

var registeredTypes *RegisteredTypes

type RegisteredTypes struct {
	Triggers       map[string]Trigger
	Apps           map[string]App
	Operators      map[string]Operator
	MagicVariables map[string]MagicVariable
	Handlers       []Handler
}

// func init() {
// 	//initialize static instance on load
// 	registeredTypes = &RegisteredTypes{
// 		Triggers:       make(map[string]Trigger),
// 		Apps:           make(map[string]App),
// 		Operators:      make(map[string]Operator),
// 		MagicVariables: make(map[string]MagicVariable),
// 	}
// }

func GetRegisteredTypes() *RegisteredTypes {
	if registeredTypes == nil {
		registeredTypes = &RegisteredTypes{
			Triggers:       make(map[string]Trigger),
			Apps:           make(map[string]App),
			Operators:      make(map[string]Operator),
			MagicVariables: make(map[string]MagicVariable),
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

func RegisterMagicVariable(magic MagicVariable) {
	name := magic.Name()
	klog.V(2).Infof("Register magic variable: %s", name)
	GetRegisteredTypes().MagicVariables[name] = magic
}

func RegisterHandler(handler Handler) {
	klog.V(2).Infof("Register handler: %s", handler.Type())
	GetRegisteredTypes().Handlers = append(GetRegisteredTypes().Handlers, handler)
}

func NewTrigger(name string, options map[string]interface{}) Trigger {
	trigger := GetRegisteredTypes().Triggers[name]
	return trigger.New(options)
}

func AppModeAvailable(app App, mode AppMode) error {
	for _, item := range app.AppModes() {
		if item == mode {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("mode %s not allowed", mode))
}

func NewApp(name string, mode AppMode, options map[string]interface{}) (App, error) {
	app := GetRegisteredTypes().Apps[name]
	if err := AppModeAvailable(app, mode); err != nil {
		return app, err
	}
	return app.New(mode, options), nil
}

func NewOperator(name string) Operator {
	operator := GetRegisteredTypes().Operators[name]
	return operator
}

func NewMagicVariable(name string, value interface{}) MagicVariable {
	name = strings.ToLower(name)
	if magic, ok := GetRegisteredTypes().MagicVariables[name]; ok {
		return magic.New(value)
	}
	return nil
}

func GetHandlers() []Handler {
	return GetRegisteredTypes().Handlers
}
