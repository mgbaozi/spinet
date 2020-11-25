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

func (SimpleInput) getExampleData() string {
	rand.Seed(time.Now().Unix())
	index := rand.Int() % 3
	contents := []string{"apple", "orange", "banana"}
	return fmt.Sprintf(`{"content": "%s"}`, contents[index])
}

func (input SimpleInput) Execute(ctx *Context, data interface{}) error {
	example := input.getExampleData()
	return json.Unmarshal([]byte(example), data)
}
