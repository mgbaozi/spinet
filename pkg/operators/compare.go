package operators

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"k8s.io/klog/v2"
	"reflect"
)

type CompareResult int

const (
	CompareResultGreater CompareResult = -1
	CompareResultEqual                 = 0
	CompareResultLess                  = 1
)

func (r CompareResult) String() string {
	switch r {
	case CompareResultGreater:
		return "greater"
	case CompareResultEqual:
		return "equal"
	case CompareResultLess:
		return "less"
	default:
		return "illegal result"
	}
}

func convertToFloat64IfPossible(value interface{}) interface{} {
	// switch value.(type) {
	// case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32:
	// 	return float64(value)

	// }
	if value == nil {
		return value
	}
	var floatType = reflect.TypeOf(float64(0))
	v := reflect.ValueOf(value)
	v = reflect.Indirect(v)
	convertible := v.Type().ConvertibleTo(floatType)
	klog.V(8).Infof("Try converting value %v to float64, result: %v", value, convertible)
	if !convertible {
		return value
	}
	fv := v.Convert(floatType)
	return fv.Float()
}

func compareString(lhs, rhs string) CompareResult {
	if lhs > rhs {
		return CompareResultGreater
	}
	if lhs < rhs {
		return CompareResultLess
	}
	return CompareResultEqual
}

func compareFloat64(lhs, rhs float64) CompareResult {
	if lhs > rhs {
		return CompareResultGreater
	}
	if lhs < rhs {
		return CompareResultLess
	}
	return CompareResultEqual
}

func compareInt(lhs, rhs int) CompareResult {
	if lhs > rhs {
		return CompareResultGreater
	}
	if lhs < rhs {
		return CompareResultLess
	}
	return CompareResultEqual
}

func compareBoolean(lhs, rhs bool) CompareResult {
	if lhs && !rhs {
		return CompareResultGreater
	}
	if !lhs && rhs {
		return CompareResultLess
	}
	return CompareResultEqual
}

func compareInRestrictedTypes(lhs, rhs interface{}) (CompareResult, error) {
	switch lhs.(type) {
	case float64:
		if rv, ok := rhs.(float64); ok {
			return compareFloat64(lhs.(float64), rv), nil
		}
	case string:
		if rv, ok := rhs.(string); ok {
			return compareString(lhs.(string), rv), nil
		}
	case bool:
		if rv, ok := rhs.(bool); ok {
			return compareBoolean(lhs.(bool), rv), nil
		}
	default:
		return CompareResultGreater, errors.New("values not have a same type")
	}
	return CompareResultGreater, errors.New("convert value to basic type failed")
}

func compareSlice(lhs, rhs interface{}) (CompareResult, error) {
	if lslice, ok := lhs.([]interface{}); ok {
		if rslice, ok := rhs.([]interface{}); ok {
			llength := len(lslice)
			rlength := len(rslice)
			if llength == 0 || rlength == 0 {
				return compareInt(llength, rlength), nil
			}
			return compare(lslice[0], rslice[0])
		}
	}
	return CompareResultGreater, errors.New("convert value to slice failed")
}

func compare(lhs, rhs interface{}) (res CompareResult, err error) {
	defer func() {
		if err != nil {
			klog.V(6).Info("Compare %v with %v failed with error: %v", lhs, rhs, err)
		} else {
			klog.V(7).Infof("Compare %v with %v, result: %s", lhs, rhs, res)
		}
	}()
	if lhs == rhs {
		return CompareResultEqual, nil
	}
	// TODO: nil can't compare with other type of value
	// TODO: but need a mechanism to process the nil value, for example fill a zero
	if lhs == nil || rhs == nil {
		if lhs == nil {
			return CompareResultLess, nil
		}
		if rhs == nil {
			return CompareResultGreater, nil
		}
	}
	lhs = convertToFloat64IfPossible(lhs)
	rhs = convertToFloat64IfPossible(rhs)
	if cmp.Equal(lhs, rhs) {
		return CompareResultEqual, nil
	}
	lvalue := reflect.ValueOf(lhs)
	rvalue := reflect.ValueOf(rhs)
	if lvalue.Type() != rvalue.Type() {
		return CompareResultGreater,
			errors.New(fmt.Sprintf("can't compare %s with %s", lvalue.Type(), rvalue.Type()))
	}
	switch lvalue.Kind() {
	case reflect.String, reflect.Float64, reflect.Bool:
		return compareInRestrictedTypes(lhs, rhs)
	case reflect.Slice:
		return compareSlice(lhs, rhs)
	default:
		return CompareResultGreater,
			errors.New(fmt.Sprintf("can't compare %s with %s", lvalue.Type(), rvalue.Type()))
	}
}

func isEqual(lhs, rhs interface{}) bool {
	if lhs == rhs {
		return true
	}
	lvalue := convertToFloat64IfPossible(lhs)
	rvalue := convertToFloat64IfPossible(rhs)
	equal := cmp.Equal(lvalue, rvalue)
	klog.V(7).Infof("Compare %v with %v, equal: %v", lvalue, rvalue, equal)
	return equal
}
