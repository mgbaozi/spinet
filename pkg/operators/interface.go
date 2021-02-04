package operators

import "k8s.io/klog/v2"

var registered map[string]Operator

func GetRegistered() map[string]Operator {
	if registered == nil {
		registered = make(map[string]Operator)
	}
	return registered
}

func register(operator Operator) {
	name := operator.Name()
	klog.V(2).Infof("Register operator: %s", name)
	GetRegistered()[name] = operator
}

func registerAll(operators ...Operator) {
	for _, operator := range operators {
		register(operator)
	}
}

type Operator interface {
	Name() string
	Do(values []interface{}) (interface{}, error)
}

func New(name string) Operator {
	operator := GetRegistered()[name]
	return operator
}
