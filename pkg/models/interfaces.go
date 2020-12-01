package models

type Input struct {
	App        App
	Conditions []Condition
}

type Output struct {
	App        App
	Conditions []Condition
}

type AppMode string

const (
	AppModeInput  AppMode = "input"
	AppModeOutPut AppMode = "output"
)

type App interface {
	Name() string
	New(options map[string]interface{}) App
	Modes() []AppMode
	Execute(mode AppMode, ctx *Context, data interface{}) error
}

type Trigger interface {
	New(options map[string]interface{}) Trigger
	Name() string
	Triggered() <-chan struct{}
}

type Operator interface {
	Name() string
	Do(values []interface{}) (bool, error)
}
