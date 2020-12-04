package models

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
)

type Condition struct {
	Operator   Operator
	Conditions []Condition
	Values     []Value
}

func (condition Condition) String() string {
	return fmt.Sprintf("Condition [op=%s,conditions=%v,values=%v]",
		condition.Operator.Name(), condition.Conditions, condition.Values)
}

func (condition Condition) Exec(dictionary, appdata interface{}) (res bool, err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("%s execute failed with error %v", condition, err)
		} else {
			klog.V(4).Infof("%s execute with result %v", condition, res)
		}
	}()
	operator := condition.Operator
	if operator == nil {
		return false, errors.New("empty operator")
	}
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
