package models

type BaseContext struct {
	Meta       Meta
	Status     Status
	Dictionary map[string]interface{}
	AppData    map[string]interface{}
	Trace      Trace
}

func NewBaseContext() *BaseContext {
	return &BaseContext{
		Dictionary: make(map[string]interface{}),
		AppData:    make(map[string]interface{}),
	}
}

func NewBaseContextWithDictionary(dictionary map[string]interface{}) *BaseContext {
	context := NewBaseContext()
	if dictionary != nil {
		context.Dictionary = dictionary
	}
	return context
}

type Shader struct{}

func NewShader() *Shader {
	return &Shader{}
}

type Context struct {
	*BaseContext
	shader Shader
}

func NewContext() Context {
	return Context{
		BaseContext: NewBaseContext(),
	}
}

func (ctx Context) Shade(shader *Shader) Context {
	if shader == nil {
		shader = NewShader()
	}
	return Context{
		ctx.BaseContext,
		*shader,
	}
}

func NewContextWithDictionary(dictionary map[string]interface{}) Context {
	context := Context{
		BaseContext: NewBaseContextWithDictionary(dictionary),
	}
	return context
}
