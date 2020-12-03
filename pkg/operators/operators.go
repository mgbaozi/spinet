package operators

import (
	"errors"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterOperators(
		EQ{},
		Contains{},
		And{},
		Or{},
	)
}

type EQ struct{}

func (EQ) Name() string {
	return "eq"
}

func (EQ) Do(values []interface{}) (bool, error) {
	klog.V(6).Infof("Compare values %v", values)
	if len(values) < 2 {
		return true, nil
	}
	for i := 0; i < len(values)-1; i++ {
		if values[i] != values[i+1] {
			return false, nil
		}
	}
	return true, nil
}

type Contains struct{}

func (Contains) Name() string {
	return "contains"
}

func listContainsValue(values []interface{}, value interface{}) bool {
	for _, item := range values {
		klog.V(6).Infof("Compare %v with %v, result: %v", item, value, item == value)
		if item == value {
			return true
		}
	}
	return false
}

func (Contains) Do(values []interface{}) (bool, error) {
	if len(values) < 2 {
		return false, errors.New("no enough arguments for operator contains")
	}
	container := values[0]
	switch container.(type) {
	case []interface{}:
		klog.V(6).Infof("Container type is list: %v", container)
		for _, value := range values[1:] {
			if !listContainsValue(container.([]interface{}), value) {
				return false, nil
			}
		}
		return true, nil
	case map[string]interface{}:
		klog.V(6).Infof("Container type is map: %v", container)
		for _, value := range values[1:] {
			if key, ok := value.(string); ok {
				if _, ok := container.(map[string]interface{})[key]; !ok {
					return false, nil
				}
			} else {
				return false, errors.New("operator contains for map need values are both string")
			}
		}
		return true, nil
	default:
		return false, errors.New("operator contains need a map or list as first argument")
	}
}

type And struct{}

func (And) Name() string {
	return "and"
}

func (And) Do(values []interface{}) (bool, error) {
	for _, value := range values {
		if res, ok := value.(bool); ok {
			if !res {
				return false, nil
			}
		} else {
			return false, errors.New("operator 'And' execute failed: can't convert value to boolean")
		}
	}
	return true, nil
}

type Or struct{}

func (Or) Name() string {
	return "or"
}

func (Or) Do(values []interface{}) (bool, error) {
	for _, value := range values {
		if res, ok := value.(bool); ok {
			if res {
				return true, nil
			}
		} else {
			return false, errors.New("operator 'Or' execute failed: can't convert value to boolean")
		}
	}
	return false, nil
}
