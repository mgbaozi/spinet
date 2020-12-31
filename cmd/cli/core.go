package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/cluster"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/urfave/cli/v2"
	"net/http"
)

func ListApps(c echo.Context) error {
	apps := models.GetApps()
	res := make([]string, 0)
	for _, app := range apps {
		res = append(res, app.AppName())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": res,
	})
}

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   err.Error(),
				"message": "handler failed with error",
			})
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
	ws.GET("/api/namespaces/:namespace/tasks/:task", cl.GetTask)
	ws.GET("/api/apps", ListApps)
	serveHTTP(ws, port)
	return nil
}
