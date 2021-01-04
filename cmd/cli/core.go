package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/cluster"
	"github.com/mgbaozi/spinet/pkg/common/rest"
	"github.com/urfave/cli/v2"
	"net/http"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			if e, ok := err.(*rest.HandlerError); ok {
				return c.JSON(e.Code, rest.ErrorResponse(e))
			}
			if e, ok := err.(*echo.HTTPError); ok {
				return c.JSON(e.Code, rest.NewResponse(e.Code, e.Message.(string), nil))
			}
			return c.JSON(http.StatusNotFound, rest.NewResponse(http.StatusNotFound, "", nil))
		}
		return err

	}
}

func core(c *cli.Context) error {
	ws := getGoEcho()
	ws.Use(ErrorHandler)
	cl := cluster.NewCluster()
	// cl.Resource.CreateNamespace("default")
	ws.GET("/api/namespaces", cl.ListNamespaces)
	ws.POST("/api/namespaces", cl.CreateNamespace)
	ws.GET("/api/namespaces/:namespace/tasks", cl.ListTasks)
	ws.POST("/api/namespaces/:namespace/tasks", cl.CreateTask)
	ws.GET("/api/namespaces/:namespace/tasks/:task", cl.GetTask)
	ws.GET("/api/apps", cluster.ListApps)
	ws.POST("/api/apps", cluster.CreateApp)
	serveHTTP(ws, port)
	return nil
}
