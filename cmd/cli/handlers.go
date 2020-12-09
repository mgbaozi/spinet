package main

import (
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/labstack/echo/v4"
	"github.com/mgbaozi/spinet/pkg/triggers"
	"k8s.io/klog/v2"
	"net/http"
)

func serveGoRest(port int) error {
	ws := new(restful.WebService)
	ws.Path("/")
	hooks := triggers.GetHookResource()
	ws.Route(
		ws.
			POST("/namespaces/{namespace}/tasks/{task}/hooks/{hook}").
			To(hooks.GoRestfulHookHandler))
	restful.Add(ws)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func serveGoEcho(port int) error {
	ws := echo.New()
	hooks := triggers.GetHookResource()
	ws.POST("/namespaces/:namespace/tasks/:task/hooks/:hook", hooks.GoEchoHookHandler)
	return ws.Start(fmt.Sprintf(":%d", port))
}

func serveHTTP() {
	// ws := getEchoService()
	klog.Fatal(serveGoEcho(8080))
	// klog.Fatal(ws.Start(":8080"))
}
