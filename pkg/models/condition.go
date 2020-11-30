package models

type Condition struct {
	Operator   Operator
	Conditions []Condition
	Values     []Value
}

func ProcessConditions(conditions []Condition, data interface{}) bool {
	return true
}

func ProcessCommonConditions(conditions []Condition, ctx *Context) bool {
	return ProcessConditions(conditions, ctx.Dictionary)
}

func ProcessAppConditions(app string, conditions []Condition, ctx *Context) bool {
	data := ctx.Data[app]
	return ProcessConditions(conditions, data)
}
