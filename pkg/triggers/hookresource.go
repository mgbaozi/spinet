package triggers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/handlers"
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
	taskName := c.Param("task")
	hookName := c.Param("hook")
	id := hookId(namespace, taskName, hookName)
	klog.V(4).Infof("Trigger hook %s", id)
	hook := h.getHook(id)
	var req, resp interface{}
	if hook == nil {
		klog.V(4).Infof("Hook %s not found", id)
		return c.JSON(http.StatusNotFound, handlers.Response{
			Meta: handlers.Meta{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("hook %s not found", id),
			},
			Data: resp,
		})
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return c.JSON(http.StatusBadRequest, handlers.Response{
			Meta: handlers.Meta{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("parse json failed with error: %v", err),
			},
			Data: resp,
		})
	}
	if err := hook.Trigger(req, &resp); err != nil {
		return c.JSON(http.StatusBadRequest, handlers.Response{
			Meta: handlers.Meta{
				Code:    http.StatusNotFound,
				Message: fmt.Sprintf("hook trigger return error: %v", err),
			},
			Data: resp,
		})
	}
	/*
		Check hook type: sync & async
		sync hook will return output data
		async hook will return this data
		consider how to give this data to task.Context
	*/
	return c.JSON(http.StatusOK, handlers.Response{
		Meta: handlers.Meta{
			Code:    http.StatusOK,
			Message: "",
		},
		Data: resp,
	})
}
