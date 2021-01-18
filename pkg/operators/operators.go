package operators

import (
	"errors"
	"github.com/mgbaozi/spinet/pkg/models"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterOperators(
		EQ{},
		Greater{},
		Less{},
		LE{},
		GE{},
		Contains{},
		And{},
		Or{},
		Not{},
	)
}

func logOperatorResult(name string, res interface{}, err error) {
	if err != nil {
		klog.V(6).Infof("Operator %s failed with error: %v", name, err)
	} else {
		klog.V(6).Infof("Operator %s success with result %v", name, res)
	}
}

type Contains struct{}

func (Contains) Name() string {
	return "contains"
}

func listContainsValue(values []interface{}, value interface{}) bool {
	for _, item := range values {
		equal := isEqual(value, item)
		if equal {
			return true
		}
	}
	return false
}

func (op Contains) Do(values []interface{}) (res interface{}, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	if len(values) < 2 {
		return false, errors.New("no enough arguments for operator contains")
	}
	container := values[0]
	switch container.(type) {
	case []interface{}:
		klog.V(7).Infof("Container type is list: %v", container)
		for _, value := range values[1:] {
			if !listContainsValue(container.([]interface{}), value) {
				return false, nil
			}
		}
		return true, nil
	case map[string]interface{}:
		klog.V(7).Infof("Container type is map: %v", container)
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

func (op And) Do(values []interface{}) (res interface{}, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
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

func (op Or) Do(values []interface{}) (res interface{}, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
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

type Not struct{}

func (Not) Name() string {
	return "not"
}

//TODO: operator not should only receive one value
func (op Not) Do(values []interface{}) (res interface{}, err error) {
	defer func() {
		logOperatorResult(op.Name(), res, err)
	}()
	for _, value := range values {
		if res, ok := value.(bool); ok {
			if res {
				return false, nil
			}
		} else {
			return false, errors.New("operator 'Not' execute failed: can't convert value to boolean")
		}
	}
	return true, nil
}
