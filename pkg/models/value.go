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
	ValueTypeConstant ValueType = "constant"
	ValueTypeVariable ValueType = "variable"
	ValueTypeTemplate ValueType = "template"
	ValueTypeMagic    ValueType = "magic"
	ValueTypeMap      ValueType = "map"
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
	case ValueTypeConstant, ValueTypeVariable, ValueTypeTemplate, ValueTypeMagic, ValueTypeMap:
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

func detectValueTypeFromString(content string) ValueType {
	if !strings.HasPrefix(content, "$") {
		return ValueTypeConstant
	}
	if strings.HasPrefix(content, "$.") {
		return ValueTypeVariable
	}
	if strings.HasPrefix(content, "${") && strings.HasSuffix(content, "}") {
		return ValueTypeTemplate
	}
	return ValueTypeMagic
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

func parseMagicVariable(content string) Value {
	name := content[1:]
	return Value{
		Type:   ValueTypeMagic,
		Source: ValueSourceNone,
		Value:  name,
	}
}

func parseTemplate(content string) Value {
	return Value{
		Type:   ValueTypeTemplate,
		Source: ValueSourceMerged,
		Value:  content[2 : len(content)-1],
	}
}

func constantValue(content interface{}) Value {
	return Value{
		Type:   ValueTypeConstant,
		Source: ValueSourceNone,
		Value:  content,
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
		case ValueTypeMagic:
			return parseMagicVariable(str)
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
				Type:   ValueTypeMap,
				Source: ValueSourceNone,
				Value:  res,
			}
		case ValueTypeVariable, ValueTypeTemplate:
			return Value{
				Type:   valueType,
				Source: detectValueSourceFromMap(dict),
				Value:  dict["value"],
			}
		default:
			return Value{
				Type:   valueType,
				Source: ValueSourceNone,
				Value:  dict["value"],
			}
		}
	}
	return constantValue(content)
}

//func ParseValue(content interface{}) Value {
//	klog.V(6).Infof("Parse value: %v", content)
//	if str, ok := content.(string); ok {
//		klog.V(7).Infof("Value type is string: %s", str)
//		if isVariable(str) {
//			return parseVariable(str)
//		}
//		if isTemplate(str) {
//			return Value{
//				Type:   ValueTypeTemplate,
//				Source: getValueSource(string(str[0])),
//				Value:  str[2 : len(str)-1],
//			}
//		}
//	}
//	if dict, ok := content.(map[string]interface{}); ok {
//		var value Value
//		if t, ok := dict["type"]; ok {
//			if value.Type, ok = toValueType(t); ok {
//				if value.Type == ValueTypeVariable || value.Type == ValueTypeTemplate {
//					if s, ok := dict["source"]; ok {
//						if value.Source, ok = toValueSource(s); !ok {
//							klog.V(2).Infof("Wrong value source %v, fallback to dictionary", value.Source)
//							value.Source = ValueSourceDictionary
//						}
//					} else {
//						value.Source = ValueSourceDictionary
//					}
//				}
//				value.Value = dict["value"]
//				return value
//			} else {
//				klog.V(8).Infof("Wrong value type %v", t)
//			}
//		}
//		// Parse as non-value map
//		res := make(map[string]interface{})
//		for key, item := range dict {
//			res[key] = ParseValue(item)
//		}
//		return Value{
//			Type:  ValueTypeMap,
//			Value: res,
//		}
//	}
//	return Value{
//		Type:   ValueTypeConstant,
//		Source: ValueSourceNone,
//		Value:  content,
//	}
//}

func (value Value) Format() interface{} {
	if value.Type == ValueTypeConstant {
		return value.Value
	}
	if value.Type == ValueTypeMap {
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

func mergeData(dictionary interface{}, appdata interface{}, super interface{}) map[string]interface{} {
	res := map[string]interface{}{
		"__dict__":  dictionary,
		"__app__":   appdata,
		"__super__": super,
	}
	merge(res, dictionary)
	merge(res, appdata)
	merge(res, super)
	return res
}

func (value Value) Extract(dictionary interface{}, appdata interface{}, super interface{}) (res interface{}, err error) {
	defer func() {
		if err != nil {
			klog.V(6).
				Infof("Extract value %v with dictionary(%v) and appdata(%v) and super(%v) failed with error: %v",
					value, dictionary, appdata, super, err)
		} else {
			klog.V(6).Infof("Extract value %v success with result %v",
				value, res)
		}
	}()
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
					if values[key], err = v.Extract(dictionary, appdata, super); err != nil {
						return values, err
					}
				} else {
					values[key] = item
				}
			}
		}
		return values, nil
	}
	var variables interface{}
	klog.V(7).Infof("Value's source is %s", value.Source)
	switch value.Source {
	case ValueSourceDictionary:
		variables = dictionary
	case ValueSourceApp:
		variables = appdata
	case ValueSourceSuper:
		variables = super
	default:
		variables = mergeData(dictionary, appdata, super)
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
