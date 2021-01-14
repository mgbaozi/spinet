package utils

import "reflect"

func SetValueToPtr(source interface{}, target interface{}) {
	val := reflect.ValueOf(target)
	if val.Kind() == reflect.Ptr {
		val.Elem().Set(reflect.ValueOf(source))
	}
}
