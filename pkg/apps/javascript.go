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
	models.RegisterApp(&Javascript{})
}

type JavascriptOptions struct {
	Script string `json:"script" yaml:"script" mapstructure:"script"`
}

type Javascript struct {
	JavascriptOptions
}

func NewJavascript(options map[string]interface{}) *Javascript {
	javascript := &Javascript{}
	if err := mapstructure.Decode(options, &javascript.JavascriptOptions); err != nil {
		klog.V(2).Infof("Merge options to javascript failed: %v", err)
	}
	return javascript
}

func (*Javascript) New(options map[string]interface{}) models.App {
	return NewJavascript(options)
}

func (*Javascript) AppName() string {
	return "javascript"
}

func (*Javascript) AppOptions() []models.AppOptionItem {
	return []models.AppOptionItem{
		{Name: "script", Type: "string", Required: true},
	}
}

func (js *Javascript) Options() map[string]interface{} {
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

func (js *Javascript) Execute(ctx models.Context, data interface{}) error {
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
