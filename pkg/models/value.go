package models

import (
	"bytes"
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
	"text/template"
)

type ValueType string
type ValueSource string

const (
	ValueTypeConstant ValueType = "constant"
	ValueTypeVariable ValueType = "variable"
	ValueTypeTemplate ValueType = "template"
)

const (
	ValueSourceNone       ValueSource = ""
	ValueSourceDictionary ValueSource = "dictionary"
	ValueSourceApp        ValueSource = "app"
)

func (vt ValueType) IsValid() bool {
	switch vt {
	case ValueTypeConstant, ValueTypeVariable, ValueTypeTemplate:
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
	Type   ValueType
	Source ValueSource
	Value  interface{}
}

func NewValue() *Value {
	return &Value{}
}

func (value *Value) Parse(content interface{}) *Value {
	newValue := ParseValue(content)
	value.Type = newValue.Type
	value.Source = newValue.Source
	value.Value = newValue.Value
	return value
}

func isVariable(content string) bool {
	// $.content || #.content
	return content == "$" || content == "#" || strings.HasPrefix(content, "$.") || strings.HasPrefix(content, "#.")
}

func isTemplate(content string) bool {
	// ${example template {{.data}} value}
	return (strings.HasPrefix(content, "${") || strings.HasPrefix(content, "#{")) &&
		strings.HasSuffix(content, "}")
}

func getValueSource(prefix string) ValueSource {
	switch prefix {
	case "#":
		return ValueSourceApp
	case "$":
		return ValueSourceDictionary
	default:
		klog.Errorf("Error when parse value with prefix: %s", prefix)
		return ValueSourceNone
	}
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
	valueSource := getValueSource(keys[0])
	if len(values) == 1 {
		return Value{
			Type:   ValueTypeVariable,
			Source: valueSource,
			Value:  values[0],
		}
	}
	return Value{
		Type:   ValueTypeVariable,
		Source: valueSource,
		Value:  values,
	}
}

func ParseValue(content interface{}) Value {
	klog.V(6).Infof("Parse value: %v", content)
	if str, ok := content.(string); ok {
		klog.V(7).Infof("Value type is string: %s", str)
		if isVariable(str) {
			return parseVariable(str)
		}
		if isTemplate(str) {
			return Value{
				Type:   ValueTypeTemplate,
				Source: getValueSource(string(str[0])),
				Value:  str[2 : len(str)-1],
			}
		}
	}
	if dict, ok := content.(map[string]interface{}); ok {
		var value Value
		if t, ok := dict["type"]; ok {
			if value.Type, ok = toValueType(t); ok {
				if value.Type == ValueTypeVariable || value.Type == ValueTypeTemplate {
					if s, ok := dict["source"]; ok {
						if value.Source, ok = toValueSource(s); !ok {
							klog.V(2).Infof("Wrong value source %v, fallback to dictionary", value.Source)
							value.Source = ValueSourceDictionary
						}
					} else {
						value.Source = ValueSourceDictionary
					}
				}
				value.Value = dict["value"]
				return value
			} else {
				klog.V(8).Infof("Wrong value type %v", t)
			}
		}
	}
	return Value{
		Type:   ValueTypeConstant,
		Source: ValueSourceNone,
		Value:  content,
	}
}

func (value Value) Format() interface{} {
	if value.Type == ValueTypeConstant {
		return value.Value
	}
	var prefix string
	switch value.Source {
	case ValueSourceDictionary:
		prefix = "$"
	case ValueSourceApp:
		prefix = "#"
	default:
		prefix = "$"
	}
	if value.Type == ValueTypeTemplate {
		// TODO: check if value is string
		return fmt.Sprintf("%s{%v}", prefix, value.Value)
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

	return fmt.Sprintf("%s.%s", prefix, format)
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
	return nil, errors.New(fmt.Sprintf("Failed convert value to string or list"))

}

func (value Value) Extract(dictionary interface{}, appdata interface{}) (res interface{}, err error) {
	defer func() {
		if err != nil {
			klog.V(6).
				Infof("Extract value %v with dictionary(%v) and appdata(%v) failed with error: %v",
					value, dictionary, appdata, err)
		} else {
			klog.V(6).Infof("Extract value %v success with result %v",
				value, res)
		}
	}()
	if value.Type == ValueTypeConstant {
		klog.V(7).Infof("Value is a constant: %v", value.Value)
		return value.Value, nil
	}
	var variables interface{}
	switch value.Source {
	case ValueSourceDictionary:
		variables = dictionary
	case ValueSourceApp:
		variables = appdata
	default:
		variables = dictionary
	}
	if value.Type == ValueTypeTemplate {
		klog.V(7).Infof("Value is a template: %v with source %s data", value.Value, value.Source)
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
	klog.V(7).Infof("Value is a variable %v with source %s data", value.Value, value.Source)
	return extractVariable(value.Value, variables)
}
