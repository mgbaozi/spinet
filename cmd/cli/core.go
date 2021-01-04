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
			return c.JSON(http.StatusBadRequest, rest.NewResponse(http.StatusBadRequest, "", nil))
		}
		return err

	}
}

func core(c *cli.Context) error {
	ws := getGoEcho()
	ws.Use(ErrorHandler)
	cl := cluster.NewCluster()
	if _, err := cl.Resource.GetNamespace("default"); err != nil {
		cl.Resource.CreateNamespace("default")
	}
	api := ws.Group("/api")
	api.GET("/namespaces", cl.ListNamespaces)
	api.POST("/namespaces", cl.CreateNamespace)
	api.GET("/apps", cluster.ListApps)
	api.POST("/apps", cluster.CreateApp)
	ns := api.Group("/namespaces")
	ns.GET("/:namespace/tasks", cl.ListTasks)
	ns.POST("/:namespace/tasks", cl.CreateTask)
	ns.GET("/:namespace/tasks/:task", cl.GetTask)
	serveHTTP(ws, port)
	return nil
}
