package cluster

import (
	"github.com/labstack/echo/v4"
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

func (cluster *Cluster) ListNamespaces(c echo.Context) error {
	namespaces := cluster.Resource.ListNamespaces()
	res := make([]string, 0)
	for _, ns := range namespaces {
		res = append(res, ns.Name)
	}
	klog.V(7).Infof("List namespaces handler return %v", res)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
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
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
}

func (cluster *Cluster) GetTask(c echo.Context) error {
	namespace := c.Param("namespace")
	name := c.Param("task")
	res, err := cluster.Resource.GetTask(name, namespace)
	if err != nil {
		return err
	}
	klog.V(7).Infof("Get task handler return %v", res)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
}

func (cluster *Cluster) CreateNamespace(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": nil,
	})
}
