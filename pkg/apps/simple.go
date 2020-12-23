package apps

import (
	"encoding/json"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
	"math/rand"
	"time"
)

func init() {
	models.RegisterApp(&Simple{})
}

type Simple struct {
	Mode    models.AppMode
	Content interface{}
}

func NewSimple(mode models.AppMode, options map[string]interface{}) *Simple {
	return &Simple{
		Mode:    mode,
		Content: options["content"],
	}
}

func (*Simple) New(mode models.AppMode, options map[string]interface{}) models.App {
	return NewSimple(mode, options)
}

func (*Simple) AppName() string {
	return "simple"
}

func (*Simple) AppModes() []models.AppMode {
	return []models.AppMode{
		models.AppModeInput,
		models.AppModeOutPut,
	}
}

func (*Simple) getExampleData() string {
	rand.Seed(time.Now().Unix())
	index := rand.Int() % 3
	contents := []string{"apple", "orange", "banana"}
	return fmt.Sprintf(`{"content": "%s"}`, contents[index])
}

func (simple *Simple) Execute(ctx *models.Context, data interface{}) error {
	if simple.Mode == models.AppModeInput {
		return simple.ExecuteAsInput(ctx, data)
	} else if simple.Mode == models.AppModeOutPut {
		return simple.ExecuteAsOutput(ctx, data)
	}
	return nil
}

func (simple *Simple) ExecuteAsInput(ctx *models.Context, data interface{}) error {
	example := simple.getExampleData()
	return json.Unmarshal([]byte(example), data)
}

func (simple *Simple) ExecuteAsOutput(ctx *models.Context, data interface{}) error {
	_, err := fmt.Println("Logging output:", simple.RenderContent(ctx.Dictionary))
	return err
}

func (simple *Simple) RenderContent(variables map[string]interface{}) interface{} {
	klog.V(4).Infof("Render content %v with variables %v", simple.Content, variables)
	content, _ := models.ParseValue(simple.Content).Extract(variables, nil, nil)
	return content
}
