package client

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/common/rest"
	"net/http"
)

type tasks struct {
	client *Client
	ns     string
}

func (client *Client) Tasks(namespace string) *tasks {
	return &tasks{
		client,
		namespace,
	}
}

func emptyResponse(pointer interface{}) (resp rest.Response) {
	resp.Data = pointer
	return
}

func (c *tasks) Create(task *apis.Task) (result *apis.Task, err error) {
	url := fmt.Sprintf("%s/api/namespaces/%s/tasks", c.client.config.server, task.Namespace)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodPost, url, nil, task, &resp)
	return
}

func (c *tasks) Get(name string) (result *apis.Task, err error) {
	url := fmt.Sprintf("%s/api/namespaces/%s/tasks/%s", c.client.config.server, c.ns, name)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodGet, url, nil, nil, &resp)
	return
}

func (c *tasks) List() (result []*apis.Task, err error) {
	url := fmt.Sprintf("%s/api/namespaces/%s/tasks", c.client.config.server, c.ns)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodGet, url, nil, nil, &resp)
	return
}
