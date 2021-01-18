package apps

import (
	"fmt"
	"github.com/mgbaozi/spinet/pkg/common/utils"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"github.com/robertkrimen/otto"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterApp(&JavaScript{})
}

type JavaScriptOptions struct {
	Script string `json:"script" yaml:"script" mapstructure:"script"`
}

type JavaScript struct {
	JavaScriptOptions
}

func NewJavaScript(options map[string]interface{}) *JavaScript {
	javascript := &JavaScript{}
	if err := mapstructure.Decode(options, &javascript.JavaScriptOptions); err != nil {
		klog.V(2).Infof("Merge options to javascript failed: %v", err)
	}
	return javascript
}

func (*JavaScript) New(options map[string]interface{}) models.App {
	return NewJavaScript(options)
}

func (*JavaScript) AppName() string {
	return "javascript"
}

func (*JavaScript) AppOptions() []models.AppOptionItem {
	return []models.AppOptionItem{
		{Name: "script", Type: "string", Required: true},
	}
}

func (js *JavaScript) Options() map[string]interface{} {
	return map[string]interface{}{
		"script": js.Script,
	}
}

func formatScript(script string) string {
	return fmt.Sprintf(`
(function() {
%s
})();
	`, script)
}

func (js *JavaScript) Execute(ctx models.Context, data interface{}) error {
	vm := otto.New()
	vm.Set("__dictionary__", ctx.Dictionary)
	if magic, ok := ctx.GetMagicVariables(); ok {
		vm.Set("__magic__", magic)
	} else {
		vm.Set("__magic__", make(map[string]interface{}))
	}
	value, err := vm.Run(formatScript(js.Script))
	if err != nil {
		klog.Error(err)
		return err
	}
	utils.SetValueToPtr(value, data)
	return nil
}
