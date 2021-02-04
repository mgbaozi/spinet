package variables

import "k8s.io/klog/v2"

var registered map[string]BuildInVariable

func GetRegistered() map[string]BuildInVariable {
	if registered == nil {
		registered = make(map[string]BuildInVariable)
	}
	return registered
}

func register(variable BuildInVariable) {
	name := variable.Name()
	klog.V(2).Infof("Register variable: %s", name)
	GetRegistered()[name] = variable
}

func registerAll(variables ...BuildInVariable) {
	for _, variable := range variables {
		register(variable)
	}
}

type BuildInVariable interface {
	New(value interface{}) BuildInVariable
	Name() string
	Data() interface{}
}

func New(name string) BuildInVariable {
	variable := GetRegistered()[name]
	return variable
}

func All() map[string]BuildInVariable {
	variable := GetRegistered()
	return variable
}

func WithOverride(override map[string]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	vars := All()
	for name, item := range vars {
		res[name] = item.New(nil).Data()
	}
	if override != nil {
		for name, item := range override {
			res[name] = item
		}
	}
	return res
}
