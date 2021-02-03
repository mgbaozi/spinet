package apps

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mgbaozi/spinet/pkg/values"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterApp(&Printer{})
}

type Printer struct {
	Content interface{}
}

func NewPrint(options map[string]interface{}) *Printer {
	return &Printer{
		Content: options["content"],
	}
}

func (*Printer) New(options map[string]interface{}) models.App {
	return NewPrint(options)
}

func (*Printer) AppName() string {
	return "print"
}

func (*Printer) AppOptions() []models.AppOptionItem {
	return []models.AppOptionItem{
		{Name: "content", Type: "any", Required: false},
	}
}

func (printer *Printer) Options() map[string]interface{} {
	return map[string]interface{}{
		"content": printer.Content,
	}
}

func (printer *Printer) Execute(ctx models.Context, data interface{}) error {
	_, err := fmt.Println("Logging output:", printer.RenderContent(ctx.MergedData()))
	return err
}

func (printer *Printer) RenderContent(variables map[string]interface{}) interface{} {
	klog.V(4).Infof("Render content %v with variables %v", printer.Content, variables)
	content, _ := values.Parse(printer.Content).Extract(variables)
	return content
}
