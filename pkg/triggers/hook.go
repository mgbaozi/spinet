package triggers

import (
	"github.com/mgbaozi/spinet/pkg/handlers"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"net/http"
	"reflect"
)

func init() {
	models.RegisterTrigger(&Hook{})
}

type HookOptions struct {
	Name   string
	Mapper models.Mapper
}

//TODO: use merge library such as mergo
func NewHookOptions(options map[string]interface{}) HookOptions {
	var mapper models.Mapper
	if mapperOptions, ok := options["mapper"]; ok {
		if m, ok := mapperOptions.(map[string]interface{}); ok {
			mapper = models.ParseMapper(m)
		}
	}
	var name string
	if nameOption, ok := options["name"]; ok {
		name = nameOption.(string)
	} else {
		name = "default"
	}
	return HookOptions{
		Name:   name,
		Mapper: mapper,
	}
}

type Hook struct {
	HookOptions
	ch      chan struct{}
	running bool
	options map[string]interface{}
	ctx     *models.Context
}

func NewHook(options map[string]interface{}) *Hook {
	return &Hook{
		HookOptions: NewHookOptions(options),
		ch:          make(chan struct{}),
		options:     options,
	}
}

func (*Hook) New(options map[string]interface{}) models.Trigger {
	return NewHook(options)
}

func (hook *Hook) Options() map[string]interface{} {
	return hook.options
}

func (hook *Hook) Meta() models.Meta {
	return hook.ctx.Meta
}

func (hook *Hook) Plural() string {
	return "hooks"
}

func (hook *Hook) Name() string {
	return hook.HookOptions.Name
}

func (hook *Hook) Methods() []string {
	return []string{http.MethodPost}
}

func (*Hook) TriggerName() string {
	return "hook"
}

func (hook *Hook) run() {
	klog.V(2).Infof("Start hook %s", hook.Name())
	handlers.Register(hook)
}

func (hook *Hook) Handler(req, resp interface{}) error {
	//FIXME: context will be replaced
	ctx := hook.ctx.Sub("handler", nil)
	ctx.SetAppData(req)
	models.ProcessMapper(ctx, hook.Mapper)
	//TODO: specify resp value
	val := reflect.ValueOf(resp)
	if val.Kind() == reflect.Ptr {
		val.Elem().Set(reflect.ValueOf(req))
	}
	hook.ch <- struct{}{}
	return nil
}

func (hook *Hook) Triggered(ctx *models.Context) <-chan struct{} {
	hook.ctx = ctx
	if !hook.running {
		hook.running = true
		hook.run()
	}
	return hook.ch
}
