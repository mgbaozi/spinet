package triggers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"k8s.io/klog/v2"
	"net/http"
)

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
