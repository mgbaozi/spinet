package apps

import (
	"context"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	"reflect"
	"time"
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
	switch metrics.Datasource {
	case "prometheus":
		address, _ := models.ParseValue(metrics.URL).Extract(ctx)
		query, _ := models.ParseValue(metrics.Query).Extract(ctx)
		client, _ := api.NewClient(api.Config{
			Address: address.(string),
		})
		v1api := v1.NewAPI(client)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		r := v1.Range{
			Start: time.Now().Add(-time.Hour),
			End:   time.Now(),
			Step:  time.Minute,
		}
		result, _, err := v1api.QueryRange(ctx, query.(string), r)
		if err != nil {
			return err
		}
		switch result.Type() {
		case model.ValMatrix:
			if matrix, ok := result.(model.Matrix); ok {
				for _, item := range matrix {
					klog.V(9).Info(item.Values)
				}
			}
		case model.ValVector:
			if vector, ok := result.(model.Vector); ok {
				klog.V(9).Info(vector)
			}
		}
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Ptr {
			val.Elem().Set(reflect.ValueOf(result))
		}
	}
	return nil
}
