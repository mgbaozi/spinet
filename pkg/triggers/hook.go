package triggers

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"net/http"
)

func init() {
	models.RegisterTrigger(&Hook{})
}

type HookResource struct {
	hooks map[string]*Hook
}

var hookResource *HookResource

func newHookResource() *HookResource {
	return &HookResource{
		hooks: make(map[string]*Hook),
	}
}

func GetHookResource() *HookResource {
	if hookResource == nil {
		hookResource = newHookResource()
	}
	return hookResource
}

func (h *HookResource) getHook(id string) *Hook {
	return h.hooks[id]
}

func (h *HookResource) Register(hook *Hook) {
	id := hook.Id()
	klog.V(2).Infof("Register hook %s", id)
	h.hooks[id] = hook
}

func (h *HookResource) Deregister(hook *Hook) {
	id := hook.Id()
	klog.V(2).Info("Deregister hook %s", id)
	h.hooks[id] = nil
}

func (h *HookResource) GoRestfulHookHandler(request *restful.Request, response *restful.Response) {
	namespace := request.PathParameter("namespace")
	task := request.PathParameter("task")
	hook := request.PathParameter("hook")
	id := hookId(namespace, task, hook)
	h.HookHandler(id, response.ResponseWriter, request.Request)
}

func (h *HookResource) GoEchoHookHandler(c echo.Context) error {

	namespace := c.Param("namespace")
	task := c.Param("task")
	hook := c.Param("hook")
	id := hookId(namespace, task, hook)
	return h.HookHandler(id, c.Response().Writer, c.Request())
}

func (h *HookResource) HookHandler(id string, w http.ResponseWriter, r *http.Request) error {
	klog.V(4).Infof("Trigger hook %s", id)
	hook := h.getHook(id)
	if hook == nil {
		w.WriteHeader(404)
		return nil
	}
	var data interface{}
	hook.Trigger(&data)
	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		_, err = w.Write([]byte(err.Error()))
		return err
	}
	_, err = w.Write(resp)
	return err
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

func (hook *Hook) Trigger(data interface{}) {
	text := `{"content": "ok"}`
	json.Unmarshal([]byte(text), data)
	hook.ch <- struct{}{}
}

func (hook *Hook) Triggered(ctx *models.Context) <-chan struct{} {
	hook.ctx = ctx
	if !hook.running {
		hook.running = true
		hook.run()
	}
	return hook.ch
}
