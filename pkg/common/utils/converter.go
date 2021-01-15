package utils

import (
	"k8s.io/klog/v2"
	"reflect"
)

func ToBool(value interface{}) bool {
	if value == nil {
		return false
	}
	rvalue := reflect.ValueOf(value)
	switch rvalue.Kind() {
	case reflect.Bool:
		return rvalue.Bool()
	case reflect.String, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return rvalue.Len() > 0
	case reflect.Ptr:
		return ToBool(rvalue.Elem().Interface())
	case reflect.Struct, reflect.Interface, reflect.UnsafePointer, reflect.Func:
		return true
	default:
		if res, ok := ConvertToFloat64IfPossible(value); ok {
			f, _ := res.(float64)
			return f != 0
		}
		return false
	}
}

func ConvertToFloat64IfPossible(value interface{}) (res interface{}, ok bool) {
	if value == nil {
		return value, false
	}
	var floatType = reflect.TypeOf(float64(0))
	v := reflect.ValueOf(value)
	v = reflect.Indirect(v)
	convertible := v.Type().ConvertibleTo(floatType)
	klog.V(8).Infof("Try converting value %v to float64, result: %v", value, convertible)
	if !convertible {
		return value, false
	}
	fv := v.Convert(floatType)
	return fv.Float(), true
}
