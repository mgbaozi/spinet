package models

import (
	"bytes"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
	"text/template"
)

type ValueType string
type ValueSource string

const (
	ValueTypeConstant   ValueType = "constant"
	ValueTypeVariable   ValueType = "variable"
	ValueTypeTemplate   ValueType = "template"
	ValueTypeBuildIn    ValueType = "buildin"
	ValueTypeMap        ValueType = "map"
	ValueTypeExpression ValueType = "expression"
)

const (
	ValueSourceNone       ValueSource = ""
	ValueSourceDictionary ValueSource = "dictionary"
	ValueSourceApp        ValueSource = "app"
	ValueSourceSuper      ValueSource = "super"
	ValueSourceMerged     ValueSource = "merged"
)

func (vt ValueType) IsValid() bool {
	switch vt {
	case ValueTypeConstant, ValueTypeVariable, ValueTypeTemplate, ValueTypeBuildIn, ValueTypeMap:
		return true
	default:
		return false
	}
}

func toValueType(t interface{}) (ValueType, bool) {
	if vt, ok := t.(string); !ok {
		return ValueTypeConstant, false
	} else {
		res := ValueType(vt)
		return res, res.IsValid()
	}
}

func (vs ValueSource) IsValid() bool {
	switch vs {
	case ValueSourceDictionary, ValueSourceApp, ValueSourceNone:
		return true
	default:
		return false
	}
}

func toValueSource(s interface{}) (ValueSource, bool) {
	if vs, ok := s.(string); !ok {
		return ValueSourceNone, false
	} else {
		res := ValueSource(vs)
		return res, res.IsValid()
	}
}

type Value struct {
	Type ValueType
	// Source ValueSource
	Value interface{}
}

func NewValue() *Value {
	return &Value{}
}

func (value *Value) Parse(content interface{}) *Value {
	newValue := ParseValue(content)
	value.Type = newValue.Type
	// value.Source = newValue.Source
	value.Value = newValue.Value
	return value
}

func detectValueTypeFromString(content string) ValueType {
	if !strings.HasPrefix(content, "$") {
		return ValueTypeConstant
	}
	if content == "$" || strings.HasPrefix(content, "$.") {
		return ValueTypeVariable
	}
	if strings.HasPrefix(content, "${") && strings.HasSuffix(content, "}") {
		return ValueTypeTemplate
	}
	return ValueTypeBuildIn
}

func detectValueTypeFromMap(content map[string]interface{}) ValueType {
	if t, ok := content["type"]; ok {
		if valueType, ok := toValueType(t); ok {
			return valueType
		}
	}
	return ValueTypeMap
}

func detectValueSourceFromMap(content map[string]interface{}) ValueSource {
	if t, ok := content["source"]; ok {
		if valueSource, ok := toValueSource(t); ok {
			return valueSource
		} else {
			klog.V(2).Infof("Wrong value source %v, fallback to merged", valueSource)
		}
	}
	return ValueSourceMerged
}

func detectValueType(content interface{}) ValueType {
	if str, ok := content.(string); ok {
		return detectValueTypeFromString(str)
	}

	if dict, ok := content.(map[string]interface{}); ok {
		return detectValueTypeFromMap(dict)
	}
	return ValueTypeConstant
}

func parseTemplate(content string) Value {
	return Value{
		Type:  ValueTypeTemplate,
		Value: content[2 : len(content)-1],
	}
}

func constantValue(content interface{}) Value {
	return Value{
		Type:  ValueTypeConstant,
		Value: content,
	}
}

func ParseValue(content interface{}) Value {
	klog.V(6).Infof("Parse value: %v", content)
	var valueType = ValueTypeConstant
	if str, ok := content.(string); ok {
		klog.V(7).Infof("Value type is string: %s", str)
		valueType = detectValueType(content)
		switch valueType {
		case ValueTypeVariable:
			return parseVariable(str)
		case ValueTypeBuildIn:
			return parseBuildInVariable(str)
		case ValueTypeTemplate:
			return parseTemplate(str)
		default:
			return constantValue(content)
		}
	}
	if dict, ok := content.(map[string]interface{}); ok {
		klog.V(7).Infof("Value type is map: %v", dict)
		valueType = detectValueType(dict)
		switch valueType {
		case ValueTypeMap:
			res := make(map[string]interface{})
			for key, item := range dict {
				res[key] = ParseValue(item)
			}
			return Value{
				Type:  ValueTypeMap,
				Value: res,
			}
		case ValueTypeVariable, ValueTypeTemplate:
			return Value{
				Type:  valueType,
				Value: dict["value"],
			}
		default:
			return Value{
				Type:  valueType,
				Value: dict["value"],
			}
		}
	}
	return constantValue(content)
}

func (value Value) Format() interface{} {
	switch value.Type {
	case ValueTypeConstant:
		return value.Value
	case ValueTypeMap:
		values := make(map[string]interface{})
		if dict, ok := value.Value.(map[string]interface{}); ok {
			for key, item := range dict {
				if v, ok := item.(Value); ok {
					values[key] = v.Format()
				} else {
					values[key] = item
				}
			}
		}
		return values
	case ValueTypeTemplate:
		return fmt.Sprintf("${%v}", value.Value)
	case ValueTypeBuildIn:
		return fmt.Sprintf("$%v", value.Value)
	case ValueTypeVariable:
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
	default:
		return value.Value
	}
}

func merge(res map[string]interface{}, item interface{}) map[string]interface{} {
	if dict, ok := item.(map[string]interface{}); ok {
		for key, value := range dict {
			res[key] = value
		}
	}
	if list, ok := item.([]interface{}); ok {
		for index, value := range list {
			key := fmt.Sprintf("%d", index)
			res[key] = value
		}
	}
	return res
}

func (value Value) Extract(variables interface{}) (res interface{}, err error) {
	defer func() {
		if err != nil {
			klog.V(6).
				Infof("Extract value %v with variables(%v) failed with error: %v",
					value, variables, err)
		} else {
			klog.V(6).Infof("Extract value %v with success with result %v",
				value, res)
		}
	}()
	if ctx, ok := variables.(Context); ok {
		variables = ctx.MergedData()
	}
	if value.Type == ValueTypeConstant {
		klog.V(7).Infof("Value is a constant: %v", value.Value)
		return value.Value, nil
	}
	if value.Type == ValueTypeMap {
		klog.V(7).Infof("Value is a map: %v", value.Value)
		values := make(map[string]interface{})
		if dict, ok := value.Value.(map[string]interface{}); ok {
			for key, item := range dict {
				if v, ok := item.(Value); ok {
					if values[key], err = v.Extract(variables); err != nil {
						return values, err
					}
				} else {
					values[key] = item
				}
			}
		}
		return values, nil
	}
	if value.Type == ValueTypeBuildIn {
		klog.V(7).Infof("Value is a build-in variable: %v", value.Value)
		if dict, ok := variables.(map[string]interface{}); ok {
			if vars, ok := dict["__buildin__"]; ok {
				return extractBuildInVariable(value.Value, vars.(map[string]interface{}))
			}
			return extractBuildInVariable(value.Value, dict)
		}
	}
	if value.Type == ValueTypeTemplate {
		klog.V(7).Infof("Value is a template: %v", value.Value)
		tmpl, err := template.New("value_parser").Parse(value.Value.(string))
		if err != nil {
			return value.Value, err
		}
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, variables)
		if err != nil {
			return value.Value, err
		}
		return buffer.String(), nil
	}
	klog.V(7).Infof("Value is a variable %v", value.Value)
	return extractVariable(value.Value, variables)
}
