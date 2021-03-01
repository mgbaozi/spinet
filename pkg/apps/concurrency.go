package apps

import (
	"github.com/mgbaozi/spinet/pkg/apis"
	"github.com/mgbaozi/spinet/pkg/common/utils"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterApp(&Concurrency{})
}

type ConcurrencyOptions struct {
	Steps interface{}
}

type Concurrency struct {
	ConcurrencyOptions
	Steps []models.Step
}

func NewConcurrency(options map[string]interface{}) *Concurrency {
	concurrency := &Concurrency{}
	if err := mapstructure.Decode(options, &concurrency.ConcurrencyOptions); err != nil {
		klog.V(2).Infof("parse options for app `concurrency` failed with error: %v", err)
	}
	yml, _ := yaml.Marshal(concurrency.ConcurrencyOptions.Steps)
	var steps []apis.Step
	if err := yaml.Unmarshal(yml, &steps); err != nil {
		klog.Errorf("Unmarshal app failed with error: %v", err)
	}
	for _, item := range steps {
		step, err := item.Parse()
		if err != nil {
			klog.Errorf("Parse step failed with error: %v", err)
		}
		concurrency.Steps = append(concurrency.Steps, step)
	}
	return concurrency
}

func (*Concurrency) New(options map[string]interface{}) models.App {
	return NewConcurrency(options)
}

func (*Concurrency) AppName() string {
	return "concurrency"
}

func (*Concurrency) AppOptions() []models.AppOptionItem {
	return models.AppOptionsFromStructPtr(&ConcurrencyOptions{})
}

func (concurrency *Concurrency) Options() (res map[string]interface{}) {
	if err := mapstructure.Decode(concurrency.ConcurrencyOptions, &res); err != nil {
		klog.Errorf("Format options for app `concurrency` failed with error: %v", err)
	}
	return
}

func (concurrency *Concurrency) Execute(ctx models.Context, data interface{}) error {
	if res, err := models.ConcurrencyProcessSteps(ctx, concurrency.Steps, nil); err != nil {
		return err
	} else {
		utils.SetValueToPtr(res, data)
	}
	return nil
}
