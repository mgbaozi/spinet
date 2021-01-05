package client

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/common/rest"
	"net/http"
)

type apps struct {
	client *Client
}

func (client *Client) Apps() *apps {
	return &apps{
		client,
	}
}

func (c *apps) Create(app *apis.CustomApp) (result *apis.CustomApp, err error) {
	url := fmt.Sprintf("%s/api/apps", c.client.config.server)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodPost, url, nil, app, &resp)
	return
}

func (c *apps) List() (result []*apis.App, err error) {
	url := fmt.Sprintf("%s/api/apps", c.client.config.server)
	resp := emptyResponse(&result)
	err = rest.Query(http.MethodGet, url, nil, nil, &resp)
	return
}
