package models

import "log"

var registeredTypes *RegisteredTypes

type RegisteredTypes struct {
	Triggers map[string]Trigger
	Apps     map[string]App
}

func init() {
	//initialize static instance on load
	registeredTypes = &RegisteredTypes{
		Triggers: make(map[string]Trigger),
		Apps:     make(map[string]App),
	}
}

func GetRegisteredTypes() *RegisteredTypes {
	return registeredTypes
}

func RegisterTrigger(trigger Trigger) {
	name := trigger.Name()
	log.Println("Register trigger:", name)
	registeredTypes.Triggers[name] = trigger
}

func RegisterApp(app App) {
	name := app.Name()
	registeredTypes.Apps[name] = app
}

func NewTrigger(name string, options map[string]interface{}) Trigger {
	log.Println("New trigger:", name)
	trigger := registeredTypes.Triggers[name]
	return trigger.New(options)
}

func NewApp(name string, options map[string]interface{}) App {
	app := registeredTypes.Apps[name]
	return app.New(options)
}
