package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/urfave/cli/v2"
	"k8s.io/klog/v2"
	"net/http"
)

type Cluster struct {
	Resource models.Resource
}

func NewCluster() *Cluster {
	return &Cluster{
		Resource: models.NewResource(),
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

func core(c *cli.Context) error {
	ws := getGoEcho()
	cluster := NewCluster()
	// cluster.Resource.CreateNamespace("default")
	ws.GET("/api/namespaces", cluster.ListNamespaces)
	serveHTTP(ws, port)
	return nil
}
