package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"strings"
)

func serveGoEcho(port int) error {
	ws := echo.New()
	ws.HideBanner = true
	ws.HidePort = true
	registeredHandlers := models.GetHandlers()
	for _, item := range registeredHandlers {
		handler := item.Handler()
		params := item.Params()
		var prefix string
		switch item.Type() {
		case models.HandlerTypeGlobal:
			prefix = ""
		case models.HandlerTypeInternal:
			prefix = "/namespaces/:namespace/tasks/:task"
		}
		path := strings.Join(append([]string{prefix}, params...), "/")
		klog.V(2).Infof("Register url path %s", path)
		for _, method := range item.Methods() {
			ws.Add(strings.ToUpper(method), path, handler)
		}
	}
	address := fmt.Sprintf(":%d", port)
	klog.V(2).Infof("http server started on %s", address)
	return ws.Start(address)
}

func serveHTTP(port int) {
	if port > 0 {
		klog.Fatal(serveGoEcho(port))
	} else {
		klog.Warning("Running without http server, hook and http output is not available...")
	}
}
