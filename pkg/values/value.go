package values

import "k8s.io/klog/v2"

/*
Data Variables:
$.__dict__: dictionary data
$.__app__: app data
$.__super__: super app data
$: merged data
${(.*)}: template (with merged data)

each.app
$.__super__.__key__ $.__super__.__index__: current index
$.__super__.__value__: value for current item
$.__super__.__collection__: hole collection
*/

type Value interface {
	New(value map[string]interface{}) Value
	Parse(str string) Value
	Type() ValueType
	Format() interface{}
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
	klog.V(7).Infof("Value %v is neither of string or map, parse as constant.", content)
	return &Constant{content}
}
