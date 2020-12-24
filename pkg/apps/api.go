package apps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"time"
)

func init() {
	models.RegisterApp(&API{})
}

type Header struct {
	Name  string
	Value interface{}
}

type API struct {
	Mode    models.AppMode
	URL     string
	Headers []Header
	Method  string
	Params  map[string]interface{}
}

func NewAPI(mode models.AppMode, options map[string]interface{}) *API {
	method, ok := options["method"].(string)
	if !ok {
		method = http.MethodGet
	}
	params, ok := options["params"].(map[string]interface{})
	if !ok {
		params = nil
	}
	var headers []Header
	optionsHeaders, ok := options["headers"].([]interface{})
	if ok {
		for _, item := range optionsHeaders {
			if h, ok := item.(map[string]interface{}); ok {
				headers = append(headers, Header{
					Name:  h["name"].(string),
					Value: h["value"],
				})
			}
		}
	}
	return &API{
		Mode:    mode,
		URL:     options["url"].(string),
		Headers: headers,
		Method:  method,
		Params:  params,
	}
}

func (*API) New(mode models.AppMode, options map[string]interface{}) models.App {
	return NewAPI(mode, options)
}

func (*API) AppName() string {
	return "api"
}

func (*API) AppModes() []models.AppMode {
	return []models.AppMode{
		models.AppModeInput,
		models.AppModeOutPut,
	}
}
func (api *API) Execute(ctx models.Context, data interface{}) (err error) {
	defer func() {
		if err != nil {
			klog.V(4).Infof("execute app %s in %s mode failed with error %v", api.AppName(), api.Mode, err)
		} else {
			klog.V(4).Infof("execute app %s in %s mode success", api.AppName(), api.Mode)
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
		value := models.ParseValue(item)
		res, err := value.Extract(ctx)
		if err != nil {
			return err
		}
		params[key] = res
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
	if method == "GET" {
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
