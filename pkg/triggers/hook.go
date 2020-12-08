package triggers

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
)

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
	klog.V(2).Infof("Start hook...")
}

func (hook *Hook) Triggered() <-chan struct{} {
	return hook.ch
}
