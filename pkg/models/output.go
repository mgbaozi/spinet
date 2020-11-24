package models

import (
	"fmt"
	"strings"
)

type Output interface {
	TriggerType() string
	Execute(ctx *Context, data interface{}) error
}

type SimpleOutput struct {
	Content string
}

func (SimpleOutput) TriggerType() string {
	return TriggerTypeActive
}

func (out SimpleOutput) Execute(ctx *Context, data interface{}) error {
	_, err := fmt.Println("Logging output:", out.RenderContent(ctx.Variables))
	return err
}

func (out SimpleOutput) RenderContent(variables map[string]string) string {
	if strings.HasPrefix(out.Content, "$.") {
		keys := strings.Split(out.Content, ".")
		for _, key := range keys[1:] {
			return variables[key]
		}
	}
	return out.Content
}
