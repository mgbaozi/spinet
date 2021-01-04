package models

import (
	"github.com/labstack/echo/v4"
)

type Step struct {
	App          App
	Mode         AppMode
	Mapper       map[string]Value
	Conditions   []Condition
	Dependencies []Step
}

type AppMode string

const (
	AppModeInput  AppMode = "input"
	AppModeOutPut AppMode = "output"
)

type App interface {
	AppName() string
	New(mode AppMode, options map[string]interface{}) App
	AppModes() []AppMode
	Options() map[string]interface{}
	Execute(ctx Context, data interface{}) error
}

type Trigger interface {
	New(options map[string]interface{}) Trigger
	TriggerName() string
	Triggered(ctx *Context) <-chan map[string]interface{}
	Options() map[string]interface{}
}

type Operator interface {
	Name() string
	Do(values []interface{}) (bool, error)
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
