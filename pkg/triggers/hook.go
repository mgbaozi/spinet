package triggers

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"reflect"
)

func init() {
	models.RegisterTrigger(&Hook{})
}

type HookOptions struct {
	Name   string
	Mapper models.Mapper
}

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
	ctx     *models.Context
}

func NewHook(options map[string]interface{}) *Hook {
	return &Hook{
		HookOptions: NewHookOptions(options),
		ch:          make(chan struct{}),
	}
}

func (*Hook) New(options map[string]interface{}) models.Trigger {
	return NewHook(options)
}

func (*Hook) Name() string {
	return "hook"
}

func (hook *Hook) run() {
	klog.V(2).Infof("Start hook %s", hook.Id())
	GetHookResource().Register(hook)
}

func hookId(namespace, task, hook string) string {
	return fmt.Sprintf("%s.%s.%s", namespace, task, hook)
}

func (hook *Hook) Id() string {
	return hookId(hook.ctx.Meta.Namespace, hook.ctx.Meta.Name, hook.HookOptions.Name)
}

func (hook *Hook) Trigger(req, resp interface{}) error {
	models.ProcessMapper(hook.ctx, hook.Mapper, req)
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
