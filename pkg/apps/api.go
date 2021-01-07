package apps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"strings"
	"time"
)

func init() {
	models.RegisterApp(&API{})
}

type Header struct {
	Name  string      `json:"name" yaml:"name" mapstructure:"name"`
	Value interface{} `json:"value" yaml:"value" mapstructure:"value"`
}

type APIOptions struct {
	URL     string                 `json:"url" yaml:"url" mapstructure:"url"`
	Headers []Header               `json:"headers" yaml:"headers" mapstructure:"headers"`
	Method  string                 `json:"method" yaml:"method" mapstructure:"method"`
	Params  map[string]interface{} `json:"params" yaml:"params" mapstructure:"params"`
}

type API struct {
	APIOptions
}

func NewAPI(options map[string]interface{}) *API {
	api := &API{}
	if err := mapstructure.Decode(options, &api.APIOptions); err != nil {
		klog.V(2).Infof("Merge options to api failed: %v", err)
	}
	if api.Method == "" {
		api.Method = http.MethodGet
	}
	api.Method = strings.ToUpper(api.Method)
	return api
}

func (*API) New(options map[string]interface{}) models.App {
	return NewAPI(options)
}

func (*API) AppName() string {
	return "api"
}

func (api *API) Options() (res map[string]interface{}) {
	if err := mapstructure.Decode(api.APIOptions, &res); err != nil {
		klog.Errorf("Format app.api's options failed with error: %v", err)
	}
	var headers []map[string]interface{}
	for _, item := range api.APIOptions.Headers {
		var header map[string]interface{}
		if err := mapstructure.Decode(item, &header); err != nil {
			klog.Errorf("Format app.api's options failed with error: %v", err)
		}
		headers = append(headers, header)
	}
	res["headers"] = headers
	return
}

func (api *API) Execute(ctx models.Context, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("execute app %s failed with error %v", api.AppName(), err)
		} else {
			klog.V(4).Infof("execute app %s success", api.AppName())
		}
	}()
	headers := make(map[string]string)
	for _, item := range api.Headers {
		value := models.ParseValue(item.Value)
		res, err := value.Extract(ctx)
		if err != nil {
			return err
		}
		headers[item.Name] = fmt.Sprint(res)
	}
	params := make(map[string]interface{})
	for key, item := range api.Params {
		var paramKey = key
		if value, err := models.ParseValue(key).Extract(ctx); err == nil {
			if pk, ok := value.(string); ok {
				paramKey = pk
			}
		}
		if value, err := models.ParseValue(item).Extract(ctx); err != nil {
			return err
		} else {
			params[paramKey] = value
		}
	}
	var method, url string
	var ok bool
	if v, err := models.ParseValue(api.Method).Extract(ctx); err != nil {
		return err
	} else {
		if method, ok = v.(string); !ok {
			method = api.Method
		}
	}
	if v, err := models.ParseValue(api.URL).Extract(ctx); err != nil {
		return err
	} else {
		if url, ok = v.(string); !ok {
			url = api.URL
		}
	}
	err = callAPI(method, url, headers, params, data)
	if err != nil {
		return err
	}
	return nil
}

func callAPI(method string, url string, headers map[string]string, params interface{}, response interface{}) error {
	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest("GET", url, nil)
	} else {
		var data []byte
		if params != nil {
			data, _ = json.Marshal(params)
			headers["content-type"] = "application/json"
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	}
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		message := fmt.Sprintf("HTTP Error: %s", resp.Status)
		return errors.New(message)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if response != nil {
		err = json.Unmarshal(body, response)
	}
	return err
}
