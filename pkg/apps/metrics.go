package apps

import (
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog/v2"
)

func init() {
	models.RegisterApp(&Metrics{})
}

type MetricsOptions struct {
	Datasource string `json:"datasource" yaml:"datasource" mapstructure:"datasource"`
	URL        string `json:"url" yaml:"url" mapstructure:"url"`
	Query      string `json:"query" yaml:"query" mapstructure:"query"`
	Duration   string `json:"duration" yaml:"duration" mapstructure:"duration"`
}

type Metrics struct {
	MetricsOptions
}

func NewMetrics(options map[string]interface{}) *Metrics {
	metrics := &Metrics{}
	if err := mapstructure.Decode(options, &metrics.MetricsOptions); err != nil {
		klog.V(2).Infof("Merge options to metrics failed: %v", err)
	}
	return metrics
}

func (*Metrics) New(options map[string]interface{}) models.App {
	return NewMetrics(options)
}

func (*Metrics) AppName() string {
	return "metrics"
}

func (*Metrics) AppOptions() []models.AppOptionItem {
	return models.AppOptionsFromStructPtr(&MetricsOptions{})
}

func (metrics *Metrics) Options() (res map[string]interface{}) {
	if err := mapstructure.Decode(metrics.MetricsOptions, &res); err != nil {
		klog.Errorf("Format app.metrics' options failed with error: %v", err)
	}
	return
}

func (metrics *Metrics) Execute(ctx models.Context, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("execute app %s failed with error %v", metrics.AppName(), err)
		} else {
			klog.V(4).Infof("execute app %s success", metrics.AppName())
		}
	}()
	//TODO: query metrics and set into data
	return nil
}
