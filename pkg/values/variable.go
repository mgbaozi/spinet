package values

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
)

type Variable struct {
	value interface{}
}

func (*Variable) New(value map[string]interface{}) Value {
	return &Variable{
		value: value["value"],
	}
}

func (*Variable) Parse(str string) Value {
	keys := strings.Split(str, ".")
	klog.V(7).Infof("Value is a variable, split keys are: %v", keys)
	var values []interface{}
	for _, key := range keys[1:] {
		if num, err := strconv.Atoi(key); err == nil {
			values = append(values, num)
		} else {
			values = append(values, key)
		}
	}
	klog.V(7).Infof("Parsed keys are: %v", values)
	if len(values) == 1 {
		return &Variable{
			value: values[0],
		}
	}
	return &Variable{
		value: values,
	}
}

func (*Variable) Type() ValueType {
	return ValueTypeVariable
}

func (variable *Variable) Format() string {
	return ""
}

func (variable *Variable) Extract(data map[string]interface{}) (interface{}, error) {
	if str, ok := variable.value.(string); ok {
		klog.V(7).Infof("Get value with key: %s", str)
		return data[str], nil
	} else if keys, ok := variable.value.([]interface{}); ok {
		klog.V(7).Infof("Get value with keys: %v", keys)
		var result interface{} = data
		for _, key := range keys {
			if str, ok := key.(string); ok {
				if res, ok := result.(map[string]interface{}); ok {
					result = res[str]
				} else {
					return nil, errors.New(fmt.Sprintf("Failed to exact key: %s", key))
				}
			} else if index, ok := key.(int); ok {
				if res, ok := result.([]interface{}); ok {
					result = res[index]
				} else {
					return nil, errors.New(fmt.Sprintf("Failed to exact index: %d", index))
				}
			}
		}
		return result, nil
	}
	return nil, errors.New(fmt.Sprintf("failed to convert value to string or list"))
}
