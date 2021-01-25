package values

import "k8s.io/klog/v2"

type Map struct {
	value interface{}
}

func (*Map) New(value map[string]interface{}) Value {
	res := make(map[string]interface{})
	for key, item := range value {
		res[key] = Parse(item)
	}
	return &Map{
		value: res,
	}
}

func (*Map) Parse(str string) Value {
	return &Map{
		value: str,
	}
}

func (*Map) Type() ValueType {
	return ValueTypeMap
}

func (variable *Map) Format() string {
	return ""
}

func (variable *Map) Extract(data map[string]interface{}) (interface{}, error) {
	klog.V(7).Infof("Value is a map: %v", variable.value)
	values := make(map[string]interface{})
	if dict, ok := variable.value.(map[string]interface{}); ok {
		for key, item := range dict {
			if v, ok := item.(Value); ok {
				var err error
				if values[key], err = v.Extract(data); err != nil {
					return values, err
				}
			} else {
				values[key] = item
			}
		}
	}
	return values, nil
}
