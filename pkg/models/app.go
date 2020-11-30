package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type AppMode string

const (
	AppModeInput  AppMode = "input"
	AppModeOutPut AppMode = "output"
)

type App interface {
	Name() string
	Modes() []AppMode
	Execute(mode AppMode, ctx *Context, data interface{}) error
}

type Simple struct {
	Content interface{}
}

func (Simple) Name() string {
	return "simple"
}

func (Simple) Modes() []AppMode {
	return []AppMode{
		AppModeInput,
		AppModeOutPut,
	}
}

func (Simple) getExampleData() string {
	rand.Seed(time.Now().Unix())
	index := rand.Int() % 3
	contents := []string{"apple", "orange", "banana"}
	return fmt.Sprintf(`{"content": "%s"}`, contents[index])
}

func (simple Simple) Execute(mode AppMode, ctx *Context, data interface{}) error {
	if mode == AppModeInput {
		return simple.ExecuteAsInput(ctx, data)
	} else if mode == AppModeOutPut {
		return simple.ExecuteAsOutput(ctx, data)
	}
	return nil
}

func (simple Simple) ExecuteAsInput(ctx *Context, data interface{}) error {
	example := simple.getExampleData()
	return json.Unmarshal([]byte(example), data)
}

func (simple Simple) ExecuteAsOutput(ctx *Context, data interface{}) error {
	_, err := fmt.Println("Logging output:", simple.RenderContent(ctx.Dictionary))
	return err
}

func (simple Simple) RenderContent(variables map[string]interface{}) interface{} {
	if content, ok := simple.Content.(string); ok {
		if strings.HasPrefix(content, "$.") {
			keys := strings.Split(content, ".")
			for _, key := range keys[1:] {
				return variables[key]
			}
		}
	}
	return simple.Content
}
