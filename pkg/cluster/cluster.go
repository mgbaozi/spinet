package cluster

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/common/rest"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"net/http"
)

type Cluster struct {
	Resource Resource
}

func NewCluster() *Cluster {
	return &Cluster{
		Resource: NewResource(),
	}
}

func jsonResponse(c echo.Context, code int, message string, data interface{}) error {
	// Define cluster error code
	return c.JSON(code, rest.NewResponse(code, message, data))
}

func (cluster *Cluster) ListNamespaces(c echo.Context) error {
	namespaces := cluster.Resource.ListNamespaces()
	res := make([]string, 0)
	for _, ns := range namespaces {
		res = append(res, ns.Name)
	}
	klog.V(7).Infof("List namespaces handler return %v", res)
	return jsonResponse(c, http.StatusOK, "", res)
}

func (cluster *Cluster) ListTasks(c echo.Context) error {
	namespace := c.Param("namespace")
	tasks, err := cluster.Resource.ListTasks(namespace)
	if err != nil {
		return err
	}
	res := make([]string, 0)
	for _, task := range tasks {
		res = append(res, task.Name)
	}
	klog.V(7).Infof("List tasks handler return %v", res)
	return jsonResponse(c, http.StatusOK, "", res)
}

func (cluster *Cluster) GetTask(c echo.Context) error {
	namespace := c.Param("namespace")
	name := c.Param("task")
	res, err := cluster.Resource.GetTask(name, namespace)
	if err != nil {
		return jsonResponse(c, http.StatusNotFound, err.Error(), nil)
	}
	klog.V(7).Infof("Get task handler return %v", res)
	return jsonResponse(c, http.StatusOK, "", res)
}

func (cluster *Cluster) CreateTask(c echo.Context) error {
	var request apis.Task
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return rest.WarpError(err, http.StatusBadRequest, "decode failed")
	}
	var task models.Task
	if task, err = request.Parse(); err != nil {
		return rest.WarpError(err, http.StatusBadRequest, "parse failed")
	}
	if err = cluster.Resource.CreateTask(&task); err != nil {
		return rest.WarpError(err, http.StatusConflict, "create failed")
	}
	return jsonResponse(c, http.StatusOK, "", request.Validate())
}

func (cluster *Cluster) CreateNamespace(c echo.Context) error {
	var request apis.Meta
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return rest.WarpError(err, http.StatusBadRequest, "decode failed")
	}
	if request.Name == "" {
		return rest.NewError(http.StatusBadRequest, "name can't be empty")
	}
	if err = cluster.Resource.CreateNamespace(request.Name); err != nil {
		return rest.WarpError(err, http.StatusConflict, "create failed")
	}
	return jsonResponse(c, http.StatusOK, "", map[string]interface{}{
		"name": request.Name,
	})
}

func CreateApp(c echo.Context) error {
	var request apis.CustomApp
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return jsonResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	var app models.CustomApp
	if app, err = request.Parse(); err != nil {
		return jsonResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	models.RegisterApp(&app)
	return jsonResponse(c, http.StatusOK, "", request.Validate())
}

func ListApps(c echo.Context) error {
	apps := models.GetApps()
	res := make([]apis.App, 0)
	for _, item := range apps {
		var app apis.App
		app.Name = item.AppName()
		app.Options = item.Options()
		app.Modes = item.AppModes()
		res = append(res, app)
	}
	return jsonResponse(c, http.StatusOK, "", res)
}
