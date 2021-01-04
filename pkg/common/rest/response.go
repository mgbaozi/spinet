package rest

import (
	"net/http"
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func DefaultResponse(code int) (resp Response) {
	resp.Meta.Code = code
	resp.Meta.Message = http.StatusText(code)
	resp.Data = map[string]interface{}{}
	return
}

func NewResponse(code int, message string, data interface{}) Response {
	if message == "" {
		message = http.StatusText(code)
	}
	if data == nil {
		data = make(map[string]interface{})
	}
	return Response{
		Meta: Meta{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}
