package models

import (
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

type MagicVariables map[string]interface{}

type BaseContext struct {
	Meta         Meta
	Status       Status
	BuildIn      map[string]interface{}
	Dictionary   map[string]interface{}
	AppData      map[string]interface{}
	MagicVarData map[string]MagicVariables
	trace        Trace
}

func NewBaseContext() *BaseContext {
	return &BaseContext{
		Dictionary:   make(map[string]interface{}),
		AppData:      make(map[string]interface{}),
		MagicVarData: make(map[string]MagicVariables),
	}
}

func NewBaseContextWithDictionary(dictionary map[string]interface{}) *BaseContext {
	context := NewBaseContext()
	if dictionary != nil {
		context.Dictionary = dictionary
	}
	context.BuildIn = buildInVariables(nil)
	return context
}

//TODO cache merged-data on each setXXX
type Shader struct {
	keys []string
}

func NewShader(key string) *Shader {
	shader := &Shader{}
	if len(key) == 0 {
		shader.keys = []string{}
	}
	shader.keys = strings.Split(key, ".")
	return shader
}

func (shader Shader) Level() int {
	return len(shader.keys)
}

func (shader Shader) Super() Shader {
	level := shader.Level()
	if level == 0 {
		return shader
	}
	return Shader{
		keys: shader.keys[:level-1],
	}
}

func (shader Shader) Sub(name string) Shader {
	keys := append(shader.keys, name)
	klog.V(7).Infof("Shader created with keys %v", keys)
	return Shader{
		keys: keys,
	}
}

func (shader Shader) Key() string {
	return strings.Join(shader.keys, ".")
}

type Context struct {
	*BaseContext
	shader Shader
}

func NewContext() Context {
	return Context{
		BaseContext: NewBaseContext(),
		shader:      Shader{},
	}
}

func (ctx Context) WithShader(shader *Shader) Context {
	if shader == nil {
		shader = NewShader("")
	}
	return Context{
		ctx.BaseContext,
		*shader,
	}
}

func (ctx Context) GetAppData() (interface{}, bool) {
	if ctx.shader.Level() == 0 {
		return nil, false
	}
	key := ctx.shader.Key()
	data, ok := ctx.AppData[key]
	return data, ok
}

func (ctx Context) SetAppData(data interface{}) {
	if ctx.shader.Level() == 0 {
		klog.Warning("Can not set app data to top level context")
		return
	}
	key := ctx.shader.Key()
	ctx.AppData[key] = data
}

func (ctx Context) GetMagicVariables() (map[string]interface{}, bool) {
	if ctx.shader.Level() == 0 {
		return nil, false
	}
	key := ctx.shader.Key()
	data, ok := ctx.MagicVarData[key]
	return data, ok
}

func (ctx Context) SetMagicVariables(data MagicVariables) {
	if ctx.shader.Level() == 0 {
		klog.Warning("Can not set app data to top level context")
		return
	}
	key := ctx.shader.Key()
	ctx.MagicVarData[key] = data
}

func (ctx Context) getTrace() *Trace {
	return &ctx.trace
}

func (ctx Context) Trace(success bool, message string, data interface{}) {
	trace := ctx.getTrace()
	trace.Push(success, message, data)
}

func (ctx Context) Super() Context {
	shader := ctx.shader.Super()
	return Context{
		ctx.BaseContext,
		shader,
	}
}

func (ctx Context) Sub(name string, magic MagicVariables) Context {
	//TODO: think if necessary to merge magic variables with super's

	// var ok bool
	// var super map[string]interface{}
	// if super, ok = ctx.GetMagicVariables(); !ok || super == nil {
	// 	super = make(map[string]interface{})
	// }
	shader := ctx.shader.Sub(name)
	context := Context{
		ctx.BaseContext,
		shader,
	}
	// for key, value := range super {
	// 	if _, ok := magic[key]; !ok {
	// 		magic[key] = value
	// 	}
	// }
	context.SetMagicVariables(magic)
	return context
}

func merge(res map[string]interface{}, item interface{}) map[string]interface{} {
	if dict, ok := item.(map[string]interface{}); ok {
		for key, value := range dict {
			res[key] = value
		}
	}
	if list, ok := item.([]interface{}); ok {
		for index, value := range list {
			key := fmt.Sprintf("%d", index)
			res[key] = value
		}
	}
	return res
}

func (ctx Context) MergedData() map[string]interface{} {
	appData, _ := ctx.GetAppData()
	magicVariables, _ := ctx.GetMagicVariables()
	res := map[string]interface{}{
		"__dict__":    ctx.Dictionary,
		"__app__":     appData,
		"__magic__":   magicVariables,
		"__buildin__": ctx.BuildIn,
	}
	merge(res, ctx.Dictionary)
	merge(res, appData)
	merge(res, magicVariables)
	return res
}

func NewContextWithDictionary(dictionary map[string]interface{}) Context {
	context := Context{
		BaseContext: NewBaseContextWithDictionary(dictionary),
	}
	return context
}

func (ctx Context) Mapper(mapper Mapper) {
	for key, value := range mapper {
		//TODO: super data
		if v, err := value.Extract(ctx.MergedData()); err == nil {
			ctx.Dictionary[key] = v
			klog.V(4).Infof("Set value %v to context.dictionary with key %s", v, key)
		}
	}
}
