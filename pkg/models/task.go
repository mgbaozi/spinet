package models

import "fmt"

type Context struct {
	Variables map[string]string
}

func (ctx *Context) GetVariable(name string) string {
	return ctx.Variables[name]
}

func (ctx *Context) SetVariable(name, value string) {
	ctx.Variables[name] = value
}

func NewContext() Context {
	return Context{
		Variables: make(map[string]string),
	}
}

type Task struct {
	Inputs []Input
	Outputs []Output
	Context  Context
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