package models

import (
	"k8s.io/klog/v2"
	"reflect"
	"strings"
)

type CustomApp struct {
	Task
	DefinedOptions []AppOptionItem
	options        map[string]Value
}

func (custom *CustomApp) New(options map[string]interface{}) App {
	app := &CustomApp{
		Task:    custom.Task,
		options: make(map[string]Value),
	}
	if app.OriginDictionary == nil {
		app.OriginDictionary = make(map[string]Value)
	}
	for key, item := range options {
		app.options[key] = ParseValue(item)
	}
	return app
}

func (custom *CustomApp) AppName() string {
	return custom.Name
}

func (custom *CustomApp) Register() {
	RegisterApp(custom)
}

func (custom *CustomApp) AppOptions() []AppOptionItem {
	return []AppOptionItem{}
}

func (custom *CustomApp) Options() map[string]interface{} {
	res := make(map[string]interface{})
	for key, value := range custom.options {
		res[key] = value.Format()
	}
	return res
}

func (custom *CustomApp) prepare(ctx Context) (err error) {
	for key, value := range custom.options {
		//TODO: super data
		custom.OriginDictionary[key] = value
	}
	dictionary := make(map[string]interface{})
	for key, item := range custom.OriginDictionary {
		dictionary[key], _ = item.Extract(ctx.MergedData())
	}
	custom.Context = NewContextWithDictionary(dictionary)
	return nil
}

func (custom *CustomApp) Execute(ctx Context, data interface{}) (err error) {
	var dict map[string]interface{}
	defer func() {
		if err != nil {
			klog.V(4).Infof("Execute app %s failed with error %v", custom.Name, err)
		}
		if dict != nil {
			val := reflect.ValueOf(data)
			if val.Kind() == reflect.Ptr {
				val.Elem().Set(reflect.ValueOf(dict))
			}
		}
		klog.V(2).Infof("App %s finished", custom.Name)
	}()
	if err := custom.prepare(ctx); err != nil {
		return err
	}
	var res bool
	if res, err = processSteps(custom.Context.Sub(string(TaskProgressInput), nil), custom.Inputs); err != nil || !res {
		return
	}
	if res, err = custom.processConditions(); err != nil || !res {
		return
	}
	if res, err = processSteps(custom.Context.Sub(string(TaskProgressOutput), nil), custom.Outputs); err != nil || !res {
		return
	}
	dict = ProcessMapper(custom.Aggregator, custom.Context.Dictionary)
	return
}

func fieldType(vfield reflect.Value) string {
	tfield := vfield.Type()
	switch tfield.Kind() {
	case reflect.Interface:
		return "any"
	case reflect.String:
		return "string"
	case reflect.Array, reflect.Slice:
		return "list"
	case reflect.Map, reflect.Struct:
		return "map"
	case reflect.Ptr:
		v := vfield.Elem()
		return fieldType(v)
	default:
		return "number"
	}
}

func AppOptionsFromStructPtr(class interface{}) (res []AppOptionItem) {
	vclass := reflect.ValueOf(class).Elem()
	tclass := vclass.Type()
	for i := 0; i < vclass.NumField(); i++ {
		var item AppOptionItem
		vfield := vclass.Field(i)
		item.Name = strings.ToLower(tclass.Field(i).Name)
		item.Type = fieldType(vfield)
		res = append(res, item)
	}
	return
}
