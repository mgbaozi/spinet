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
		return jsonResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	var task models.Task
	if task, err = request.Parse(); err != nil {
		return jsonResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if err = cluster.Resource.CreateTask(&task); err != nil {
		return jsonResponse(c, http.StatusConflict, err.Error(), nil)
	}
	return jsonResponse(c, http.StatusOK, "", request.Validate())
}

func (cluster *Cluster) CreateNamespace(c echo.Context) error {
	var request apis.Meta
	err := json.NewDecoder(c.Request().Body).Decode(&request)
	if err != nil {
		return jsonResponse(c, http.StatusBadRequest, err.Error(), nil)
	}
	if request.Name == "" {
		return jsonResponse(c, http.StatusBadRequest, "name can't be empty", nil)
	}
	if err = cluster.Resource.CreateNamespace(request.Name); err != nil {
		return jsonResponse(c, http.StatusConflict, err.Error(), nil)
	}
	return jsonResponse(c, http.StatusOK, "", map[string]interface{}{
		"name": request.Name,
	})
}
