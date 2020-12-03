package models

import (
	"errors"
	"k8s.io/klog/v2"
)

type Condition struct {
	Operator   Operator
	Conditions []Condition
	Values     []Value
}

func (condition Condition) Exec(dictionary, appdata interface{}) (bool, error) {
	operator := condition.Operator
	if operator == nil {
		return false, errors.New("empty operator")
	}
	klog.V(4).Infof("Process condition with operator %s and values %v", operator.Name(), condition.Values)
	if len(condition.Conditions) > 0 {
		return ProcessConditions(operator, condition.Conditions, dictionary, appdata)
	}
	var values []interface{}
	for _, value := range condition.Values {
		extracted, err := value.Extract(dictionary, appdata)
		if err != nil {
			return false, err
		}
		values = append(values, extracted)
	}
	return operator.Do(values)
}

func ProcessConditions(operator Operator, conditions []Condition, dictionary, appdata interface{}) (bool, error) {
	var values []interface{}
	for _, condition := range conditions {
		res, err := condition.Exec(dictionary, appdata)
		if err != nil {
			return false, err
		}
		values = append(values, res)
	}
	return operator.Do(values)
}
