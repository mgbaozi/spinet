package models

import (
	"github.com/labstack/echo/v4"
)

type Input struct {
	App        App
	Mapper     map[string]Value
	Conditions []Condition
}

type Output struct {
	App        App
	Mapper     map[string]Value
	Conditions []Condition
}

type AppMode string

const (
	AppModeInput  AppMode = "input"
	AppModeOutPut AppMode = "output"
)

type App interface {
	AppName() string
	New(options map[string]interface{}) App
	Modes() []AppMode
	Execute(mode AppMode, ctx *Context, data interface{}) error
}

type Trigger interface {
	New(options map[string]interface{}) Trigger
	TriggerName() string
	Triggered(ctx *Context) <-chan struct{}
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
	// TODO: use func (params ...interface{}) error
	Handler() func(c echo.Context) error
}
