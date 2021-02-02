package models

import (
	"github.com/mgbaozi/spinet/pkg/values"
	"k8s.io/klog/v2"
)

type Mapper map[string]values.Value

func ParseMapper(data map[string]interface{}) Mapper {
	mapper := make(Mapper)
	for key, content := range data {
		value := values.Parse(content)
		mapper[key] = value
	}
	return mapper
}

func ProcessMapper(mapper Mapper, variables interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range mapper {
		//FIXME: force cast will cause crash
		if v, err := value.Extract(variables.(map[string]interface{})); err == nil {
			res[key] = v
			klog.V(4).Infof("Mapper set value %v to key %s", v, key)
		}
	}
	return res
}
