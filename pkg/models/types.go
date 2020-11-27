package models

type Context struct {
	Dictionary map[string]interface{}
}

type Task struct {
	Name string
	Inputs []Input
	Outputs []Output
	Context Context
}

type Trigger interface {
	Triggered() <-chan struct{}
}

type App interface {
	Options() map[string]interface{}
	SetOptions(options map[string]interface{})
}

type Value struct {
	Type string
	Value interface{}
}

type Condition struct {
	Operator string
	Type string
	Conditions []Condition
}

type Comparator interface {
	Do(ctx *Context, values ...Value) bool
}