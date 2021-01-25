package values

import (
	"bytes"
	"k8s.io/klog/v2"
	"text/template"
)

type Template struct {
	value interface{}
}

func (*Template) New(value map[string]interface{}) Value {
	return &Template{
		value: value["value"],
	}
}

func (*Template) Parse(str string) Value {
	return &Template{
		value: str[2 : len(str)-1],
	}
}

func (*Template) Type() ValueType {
	return ValueTypeTemplate
}

func (variable *Template) Format() string {
	return ""
}

func (variable *Template) Extract(data map[string]interface{}) (interface{}, error) {
	klog.V(7).Infof("Value is a template: %v", variable.value)
	tmpl, err := template.New("value_parser").Parse(variable.value.(string))
	if err != nil {
		return variable.value, err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		return variable.value, err
	}
	return buffer.String(), nil
}
