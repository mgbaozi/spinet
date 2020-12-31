package handlers

import "github.com/mgbaozi/spinet/pkg/models"

type CustomHandler interface {
	Meta() models.Meta
	Name() string
	Plural() string
	Methods() []string
	Handler(req, resp interface{}) error
}
