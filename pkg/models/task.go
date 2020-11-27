package models

import "fmt"

func (ctx *Context) GetVariable(name string) interface{} {
	return ctx.Dictionary[name]
}

func (ctx *Context) SetVariable(name, value string) {
	ctx.Dictionary[name] = value
}

func NewContext() Context {
	return Context{
		Dictionary: make(map[string]interface{}),
	}
}

func (task *Task) Execute() {
	for _, input := range task.Inputs {
		switch input.TriggerType() {
		case TriggerTypeActive:
			fmt.Println("Active trigger")
			var data map[string]interface{}
			err := input.Execute(&task.Context, &data)
			if err != nil {
				fmt.Println(err)
			}
		case TriggerTypePassive:
			fmt.Println("Passive trigger")
		}
	}
	for _, output := range task.Outputs {
		switch output.TriggerType() {
		case TriggerTypeActive:
			fmt.Println("Active trigger")
			err := output.Execute(&task.Context, nil)
			if err != nil {
				fmt.Println(err)
			}
		case TriggerTypePassive:
			fmt.Println("Passive trigger")
		}
	}
}