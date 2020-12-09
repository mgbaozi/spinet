package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/triggers"
	"k8s.io/klog/v2"
)

func serveGoEcho(port int) error {
	ws := echo.New()
	ws.HideBanner = true
	ws.HidePort = true
	hooks := triggers.GetHookResource()
	ws.POST("/namespaces/:namespace/tasks/:task/hooks/:hook", hooks.GoEchoHookHandler)
	address := fmt.Sprintf(":%d", port)
	klog.V(2).Infof("http server started on %s", address)
	return ws.Start(address)
}

func serveHTTP(port int) error {
	return serveGoEcho(8080)
}
