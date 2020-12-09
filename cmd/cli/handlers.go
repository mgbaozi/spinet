package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/triggers"
	"k8s.io/klog/v2"
)

func serveGoEcho(port int) error {
	ws := echo.New()
	hooks := triggers.GetHookResource()
	ws.POST("/namespaces/:namespace/tasks/:task/hooks/:hook", hooks.GoEchoHookHandler)
	return ws.Start(fmt.Sprintf(":%d", port))
}

func serveHTTP() {
	klog.Fatal(serveGoEcho(8080))
}
