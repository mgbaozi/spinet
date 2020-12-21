package apps

import (
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/models"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
	"reflect"
)

func init() {
	models.RegisterApp(&Each{})
}

type Each struct {
	Mode       models.AppMode
	Collection interface{}
	Steps      []models.Step
}

func NewEach(mode models.AppMode, options map[string]interface{}) *Each {
	each := &Each{
		Mode:       mode,
		Collection: options["collection"].(string),
	}
	yml, _ := yaml.Marshal(options["apps"])
	var steps []apis.Step
	if err := yaml.Unmarshal(yml, &steps); err != nil {
		klog.Errorf("Unmarshal app failed with error: %v", err)
	}
	for _, item := range steps {
		step, err := item.Parse(mode)
		if err != nil {
			klog.Errorf("Parse step failed with error: %v", err)
		}
		each.Steps = append(each.Steps, step)
	}
	return each
}

func (*Each) New(mode models.AppMode, options map[string]interface{}) models.App {
	return NewEach(mode, options)
}

func (*Each) AppName() string {
	return "each"
}

func (*Each) AppModes() []models.AppMode {
	return []models.AppMode{
		models.AppModeInput,
		models.AppModeOutPut,
	}
}

func (each *Each) executeApps(ctx *models.Context, key interface{}, value interface{}) (results []interface{}, err error) {
	klog.V(6).Infof("Execute app with key: %v, value: %v", key, value)
	for _, step := range each.Steps {
		var res bool
		//TODO: send key & value to app
		if res, err = step.Process(ctx, fmt.Sprintf(".%d", key)); err != nil {
			return
		} else {
			results = append(results, res)
		}
	}
	return
}

func (each *Each) Execute(ctx *models.Context, data interface{}) error {
	collection, err := models.ParseValue(each.Collection).Extract(ctx.Dictionary, nil)
	if err != nil {
		return err
	}
	if l, ok := collection.([]interface{}); ok {
		klog.V(6).Infof("Collection is a list: %v", l)
		for index, item := range l {
			if res, err := each.executeApps(ctx, index, item); err != nil {
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
			if res, err := each.executeApps(ctx, key, item); err != nil {
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
