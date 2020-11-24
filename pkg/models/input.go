package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type Input interface {
	TriggerType() string
	Execute(ctx *Context, data interface{}) error
}

type SimpleInput struct {
}

func (SimpleInput) TriggerType() string{
	return TriggerTypeActive
}

func (SimpleInput) Execute(ctx *Context, data interface{}) error {
	rand.Seed(time.Now().Unix())
	contents := []string{"apple", "orange", "banana"}
	index := rand.Int() % 3
	example := fmt.Sprintf(`{"content": "%s"}`, contents[index])
	ctx.Variables["content"] = contents[index]
	return json.Unmarshal([]byte(example), data)
}
