package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"net/http"
)

func init() {
	models.RegisterHandler(GetResource())
}

var resource *Resource

type RestHandler interface {
	Context() *models.Context
	Name() string
	Plural() string
	Methods() []string
	Handler(req, resp interface{}) error
}

func resourceKey(namespace, task, plural, name string) string {
	return fmt.Sprintf("%s.%s.%s.%s", namespace, task, plural, name)
}

func restHandlerKey(handler RestHandler) string {
	ctx := handler.Context()
	return resourceKey(ctx.Meta.Namespace, ctx.Meta.Name, handler.Plural(), handler.Name())
}

type Resource struct {
	handlers map[string]RestHandler
}

func newResource() *Resource {
	return &Resource{
		handlers: make(map[string]RestHandler),
	}
}

func GetResource() *Resource {
	if resource == nil {
		resource = newResource()
	}
	return resource
}

func Register(handler RestHandler) {
	resource := GetResource()
	resource.Register(handler)
}

func (r *Resource) Register(handler RestHandler) {
	key := restHandlerKey(handler)
	klog.V(2).Infof("Register rest handler %s", key)
	r.handlers[key] = handler
}

func (r *Resource) Handler() func(c echo.Context) error {
	return r.RestHandler
}

func (r *Resource) Type() models.HandlerType {
	return models.HandlerTypeInternal
}

func (r *Resource) Methods() []string {
	return []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}
}

func (r *Resource) Params() []string {
	return []string{":plural", ":name"}
}

func methodAllowed(method string, handler RestHandler) bool {
	for _, item := range handler.Methods() {
		if method == item {
			return true
		}
	}
	return false
}

func jsonResponse(c echo.Context, code int, message string, data interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}
	return c.JSON(code, Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	})
}

func (r *Resource) RestHandler(c echo.Context) error {
	namespace := c.Param("namespace")
	task := c.Param("task")
	plural := c.Param("plural")
	name := c.Param("name")
	key := resourceKey(namespace, task, plural, name)
	klog.V(4).Infof("RestHandler for key %s", key)
	handler, ok := r.handlers[key]
	var req, resp interface{}
	if !ok {
		klog.V(4).Infof("Handler %s not found", key)
		return jsonResponse(c, http.StatusNotFound, fmt.Sprintf("Handler %s not found", key), resp)
	}
	if !methodAllowed(c.Request().Method, handler) {
		return jsonResponse(c, http.StatusMethodNotAllowed,
			fmt.Sprintf("Method %s not allowed for handler %s", c.Request().Method, key),
			resp)
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		return jsonResponse(c, http.StatusBadRequest,
			fmt.Sprintf("parse json failed with error: %v", err),
			resp)
	}
	if err := handler.Handler(req, &resp); err != nil {
		return jsonResponse(c, http.StatusBadRequest,
			fmt.Sprintf("%s handler return error: %v", plural, err),
			resp)
	}
	/*
		Check hook type: sync & async
		sync hook will return output data
		async hook will return this data
		consider how to give this data to task.Context
	*/
	return jsonResponse(c, http.StatusOK, "", resp)
}
