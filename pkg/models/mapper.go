package models

import "k8s.io/klog/v2"

type Mapper map[string]Value

func ParseMapper(data map[string]interface{}) Mapper {
	mapper := make(Mapper)
	for key, content := range data {
		value := ParseValue(content)
		mapper[key] = value
	}
	return mapper
}

func ProcessMapper(mapper Mapper, variables interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range mapper {
		if v, err := value.Extract(variables); err == nil {
			res[key] = v
			klog.V(2).Infof("Mapper set value %v to key %s", v, key)
		}
	}
	return res
}
