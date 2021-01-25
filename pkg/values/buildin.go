package values

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"strings"
)

type BuildIn struct {
	value interface{}
}

func (*BuildIn) New(value map[string]interface{}) Value {
	return &BuildIn{
		value: value["value"],
	}
}

func (*BuildIn) Parse(str string) Value {
	name := str[1:]
	return &BuildIn{
		value: name,
	}
}

func (*BuildIn) Type() ValueType {
	return ValueTypeBuildIn
}

func (variable *BuildIn) Format() string {
	return ""
}

func (variable *BuildIn) extract(variables map[string]interface{}) (res interface{}, err error) {
	if str, ok := variable.value.(string); ok {
		key := strings.ToLower(str)
		if v, ok := variables[key]; ok {
			return v, nil
		}
		return variable.value, errors.New(fmt.Sprintf("build-in variable %s not found", str))
	}
	return nil, errors.New("wrong build-in variable format")
}

func (variable *BuildIn) Extract(data map[string]interface{}) (interface{}, error) {
	klog.V(7).Infof("Value is a build-in variable: %v", variable.value)
	if vars, ok := data["__buildin__"]; ok {
		return variable.extract(vars.(map[string]interface{}))
	}
	return variable.extract(data)
}
