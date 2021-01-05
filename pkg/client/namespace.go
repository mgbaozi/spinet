package client

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/common/rest"
	"net/http"
)

type namespaces struct {
	client *Client
}

func (client *Client) Namespaces() *namespaces {
	return &namespaces{
		client,
	}
}

func (c *namespaces) Create(namespace *apis.Namespace) (result apis.Namespace, err error) {
	url := fmt.Sprintf("%s/api/namespaces", c.client.config.server)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodPost, url, nil, namespace, &resp)
	return
}

func (c *namespaces) Get(name string) (result apis.Namespace, err error) {
	url := fmt.Sprintf("%s/api/namespaces/%s", c.client.config.server, name)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodGet, url, nil, nil, &resp)
	return
}

func (c *namespaces) List() (result []*apis.Namespace, err error) {
	url := fmt.Sprintf("%s/api/namespaces", c.client.config.server)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodGet, url, nil, nil, &resp)
	return
}
