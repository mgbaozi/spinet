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

func (c *tasks) Create(task *apis.Task) (result *apis.Task, err error) {
	task.Validate()
	url := fmt.Sprintf("%s/api/namespaces/%s/tasks", c.client.config.server, task.Namespace)
	err = rest.Query(http.MethodPost, url, nil, task, result)
	return
}
