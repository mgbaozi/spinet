package models

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
)

/*
Data Variables:
$.__dict__: dictionary data
$.__app__: app data
$.__super__: super app data
$: merged data
${(.*)}: template (with merged data)

each.app
$.__super__.__key__ $.__super__.__index__: current index
$.__super__.__value__: value for current item
$.__super__.__collection__: hole collection
*/

func buildInVariables(override map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	vars := GetBuildInVariables()
	for name, item := range vars {
		res[name] = item.New(nil).Data()
	}
	if override != nil {
		for name, item := range override {
			res[name] = item
		}
	}
	return res
}

func parseBuildInVariable(content string) Value {
	name := content[1:]
	return Value{
		Type:  ValueTypeBuildIn,
		Value: name,
	}
}

func extractBuildInVariable(value interface{}, variables map[string]interface{}) (res interface{}, err error) {
	if str, ok := value.(string); ok {
		if v, ok := variables[str]; ok {
			return v, nil
		}
		return value, errors.New(fmt.Sprintf("build-in variable %s not found", str))
	}
	return nil, errors.New("wrong build-in variable format")
}

func parseVariable(str string) Value {
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
		return Value{
			Type:  ValueTypeVariable,
			Value: values[0],
		}
	}
	return Value{
		Type:  ValueTypeVariable,
		Value: values,
	}
}

func extractVariable(value interface{}, variables interface{}) (res interface{}, err error) {
	if str, ok := value.(string); ok {
		klog.V(7).Infof("Get value with key: %s", str)
		return variables.(map[string]interface{})[str], nil
	} else if keys, ok := value.([]interface{}); ok {
		klog.V(7).Infof("Get value with keys: %v", keys)
		var result = variables
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
