package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ValueType string

const (
	ValueTypeConstant ValueType = "constant"
	ValueTypeVariable ValueType = "variable"
)

type Value struct {
	Type  ValueType
	Value interface{}
}

func ParseValue(content interface{}) Value {
	if str, ok := content.(string); ok {
		if strings.HasPrefix(str, "$.") {
			keys := strings.Split(str, ".")
			var values []interface{}
			for _, key := range keys {
				if num, err := strconv.Atoi(key); err != nil {
					values = append(values, num)
				} else {
					values = append(values, key)
				}
			}
			if len(values) == 1 {
				return Value{
					Type:  ValueTypeVariable,
					Value: values[0],
				}
			} else {
				return Value{
					Type:  ValueTypeVariable,
					Value: values,
				}
			}
		}
	}
	return Value{
		Type:  ValueTypeConstant,
		Value: content,
	}
}

func (value Value) Format() interface{} {
	if value.Type == ValueTypeConstant {
		return value.Value
	}
	var format string
	if str, ok := value.Value.(string); ok {
		format = str
	} else if keys, ok := value.Value.([]interface{}); ok {
		var values []string
		for _, key := range keys {
			if str, ok := key.(string); ok {
				values = append(values, str)
			} else if num, ok := key.(int); ok {
				values = append(values, strconv.Itoa(num))
			}
		}
		format = strings.Join(values, ".")
	}
	return fmt.Sprintf("$.%s", format)
}

func (value Value) Equals(right Value) bool {
	return true
}

func (value Value) Extract(variables map[string]interface{}) (interface{}, error) {
	if value.Type == ValueTypeConstant {
		return value.Value, nil
	}
	if str, ok := value.Value.(string); ok {
		return variables[str], nil
	} else if keys, ok := value.Value.([]interface{}); ok {
		var result interface{} = variables
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
	return nil, errors.New(fmt.Sprintf("Failed convert value to string or list"))
}
