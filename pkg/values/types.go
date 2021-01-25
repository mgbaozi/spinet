package values

import (
	"strings"
)

type ValueType string

const (
	ValueTypeConstant   ValueType = "constant"
	ValueTypeVariable   ValueType = "variable"
	ValueTypeTemplate   ValueType = "template"
	ValueTypeBuildIn    ValueType = "buildin"
	ValueTypeMap        ValueType = "map"
	ValueTypeExpression ValueType = "expression"
)

func produceValue(valueType ValueType) Value {
	switch valueType {
	case ValueTypeConstant:
		return &Constant{}
	case ValueTypeVariable:
		return &Variable{}
	case ValueTypeTemplate:
		return &Template{}
	case ValueTypeBuildIn:
		return &BuildIn{}
	case ValueTypeMap:
		return &Map{}
	case ValueTypeExpression:
		return &Expression{}
	default:
		return &Constant{}
	}
}

func (vt ValueType) IsValid() bool {
	switch vt {
	case ValueTypeConstant, ValueTypeVariable, ValueTypeTemplate, ValueTypeBuildIn, ValueTypeMap, ValueTypeExpression:
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

func detectValueType(content interface{}) ValueType {
	if str, ok := content.(string); ok {
		return detectValueTypeFromString(str)
	}

	if dict, ok := content.(map[string]interface{}); ok {
		return detectValueTypeFromMap(dict)
	}
	return ValueTypeConstant
}
