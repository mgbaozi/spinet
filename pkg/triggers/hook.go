package triggers

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
)

type HookResource struct {
	hooks map[string]*Hook
}

func (h HookResource) getHook(name string) *Hook {
	return h.hooks[name]
}

func (h HookResource) HookHandler(request *restful.Request, response *restful.Response) {
	name := request.PathParameter("hook")
	hook := h.getHook(name)
	var data interface{}
	hook.Trigger(&data)
	/*
		Check hook type: sync & async
		sync hook will return output data
		async hook will return this data
		consider how to give this data to task.Context
	*/
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

func (hook *Hook) Trigger(data interface{}) {
	hook.ch <- struct{}{}
}

func (hook *Hook) Triggered(ctx *models.Context) <-chan struct{} {
	return hook.ch
}
