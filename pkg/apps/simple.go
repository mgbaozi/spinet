package apps

import (
	"encoding/json"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"math/rand"
	"time"
)

func init() {
	models.RegisterApp(&Simple{})
}

type Simple struct {
	Content interface{}
}

func NewSimple(options map[string]interface{}) *Simple {
	return &Simple{
		Content: options["content"],
	}
}

func (*Simple) New(options map[string]interface{}) models.App {
	return NewSimple(options)
}

func (*Simple) AppName() string {
	return "simple"
}

func (simple *Simple) Options() map[string]interface{} {
	return map[string]interface{}{
		"content": simple.Content,
	}
}

func (*Simple) getExampleData() string {
	rand.Seed(time.Now().Unix())
	index := rand.Int() % 3
	contents := []string{"apple", "orange", "banana"}
	return fmt.Sprintf(`{"content": "%s"}`, contents[index])
}

func (simple *Simple) Execute(ctx models.Context, data interface{}) error {
	example := simple.getExampleData()
	return json.Unmarshal([]byte(example), data)
}
