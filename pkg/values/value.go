package values

import "k8s.io/klog/v2"

type Value interface {
	New(value map[string]interface{}) Value
	Parse(str string) Value
	Type() ValueType
	Format() string
	Extract(data map[string]interface{}) (interface{}, error)
}

func Parse(content interface{}) Value {
	klog.V(6).Infof("Parse value: %v", content)
	valueType := detectValueType(content)
	klog.V(7).Infof("Value is a %s", valueType)
	value := produceValue(valueType)
	if str, ok := content.(string); ok {
		klog.V(7).Infof("Value type is string: %s", str)
		return value.Parse(str)
	}
	if dict, ok := content.(map[string]interface{}); ok {
		klog.V(7).Infof("Value type is map: %v", dict)
		return value.New(dict)
	}
	klog.Errorf("Value %v is neither of string or map, parse failed.", content)
	return nil
}
