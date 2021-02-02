package models

import (
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/values"
)

type Step struct {
	App          App
	Mapper       map[string]values.Value
	Conditions   []Condition
	Dependencies []Step
}

type TaskProgress string

const (
	TaskProgressInput      TaskProgress = "input"
	TaskProgressCondition  TaskProgress = "condition"
	TaskProgressOutput     TaskProgress = "output"
	TaskProgressAggregator TaskProgress = "aggregator"
)

type App interface {
	AppName() string
	AppOptions() []AppOptionItem
	New(options map[string]interface{}) App
	Options() map[string]interface{}
	Execute(ctx Context, data interface{}) error
}

type AppOptionItem struct {
	Name     string
	Type     string
	Required bool
}

type Trigger interface {
	New(options map[string]interface{}) Trigger
	TriggerName() string
	Triggered(ctx *Context) <-chan map[string]interface{}
	Options() map[string]interface{}
}

type Operator interface {
	Name() string
	Do(values []interface{}) (interface{}, error)
}

type HandlerType string

const (
	HandlerTypeInternal HandlerType = "internal"
	HandlerTypeGlobal   HandlerType = "global"
)

type Handler interface {
	Type() HandlerType
	Methods() []string
	Params() []string
	Handler() func(c echo.Context) error
}

type BuildInVariable interface {
	New(value interface{}) BuildInVariable
	Name() string
	Data() interface{}
}
