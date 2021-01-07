package apps

import (
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
	"reflect"
)

func init() {
	models.RegisterApp(&Each{})
}

type EachOptions struct {
	Collection interface{}
	Apps       interface{}
}

type Each struct {
	EachOptions
	Steps []models.Step
}

func NewEach(options map[string]interface{}) *Each {
	each := &Each{}
	if err := mapstructure.Decode(options, &each.EachOptions); err != nil {
		klog.V(2).Infof("parse options for app `each` failed with error: %v", err)
	}
	yml, _ := yaml.Marshal(each.EachOptions.Apps)
	var steps []apis.Step
	if err := yaml.Unmarshal(yml, &steps); err != nil {
		klog.Errorf("Unmarshal app failed with error: %v", err)
	}
	for _, item := range steps {
		step, err := item.Parse()
		if err != nil {
			klog.Errorf("Parse step failed with error: %v", err)
		}
		each.Steps = append(each.Steps, step)
	}
	return each
}

func (*Each) New(options map[string]interface{}) models.App {
	return NewEach(options)
}

func (*Each) AppName() string {
	return "each"
}

func (each *Each) Options() (res map[string]interface{}) {
	if err := mapstructure.Decode(each.EachOptions, &res); err != nil {
		klog.Errorf("Format options for app `each` failed with error: %v", err)
	}
	return
}

func (each *Each) executeApps(ctx models.Context, key interface{}, value interface{}, collection interface{}) (results []interface{}, err error) {
	klog.V(6).Infof("Execute app with key: %v, value: %v", key, value)
	for _, step := range each.Steps {
		var res bool
		name := fmt.Sprintf("each-%v", key)
		magic := map[string]interface{}{
			"__index__":      key,
			"__key__":        key,
			"__value__":      value,
			"__collection__": collection,
		}
		if res, err = step.Process(ctx.Sub(name, magic)); err != nil {
			return
		} else {
			results = append(results, res)
		}
	}
	return
}

func (each *Each) Execute(ctx models.Context, data interface{}) error {
	collection, err := models.ParseValue(each.Collection).Extract(ctx)
	if err != nil {
		return err
	}
	if l, ok := collection.([]interface{}); ok {
		klog.V(6).Infof("Collection is a list: %v", l)
		for index, item := range l {
			if res, err := each.executeApps(ctx, index, item, collection); err != nil {
				return err
			} else {
				val := reflect.ValueOf(data)
				if val.Kind() == reflect.Ptr {
					val.Elem().Set(reflect.ValueOf(res))
				}
			}
		}
		return nil
	}
	if m, ok := collection.(map[string]interface{}); ok {
		klog.V(6).Infof("Collection is a map: %v", m)
		for key, item := range m {
			if res, err := each.executeApps(ctx, key, item, collection); err != nil {
				return err
			} else {
				val := reflect.ValueOf(data)
				if val.Kind() == reflect.Ptr {
					val.Elem().Set(reflect.ValueOf(res))
				}
			}
		}
		return nil
	}
	klog.V(6).Infof("Collection is not a list: %v", collection)
	return errors.New("not a collection")
}
