package apps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mgbaozi/spinet/pkg/models"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {
	models.RegisterApp(&API{})
}

type Header struct {
	Name  string
	Value string
}

type API struct {
	URL     string
	Headers []Header
	Method  string
	Params  map[string]interface{}
}

func NewAPI(options map[string]interface{}) *API {
	method, ok := options["method"].(string)
	if !ok {
		method = http.MethodGet
	}
	return &API{
		URL:     options["url"].(string),
		Headers: nil,
		Method:  method,
	}
}

func (*API) New(options map[string]interface{}) models.App {
	return NewAPI(options)
}

func (*API) Name() string {
	return "api"
}

func (*API) Modes() []models.AppMode {
	return []models.AppMode{
		models.AppModeInput,
		models.AppModeOutPut,
	}
}

func (api *API) Execute(mode models.AppMode, ctx *models.Context, data interface{}) error {
	headers := map[string]string{
		"Authorization": "Token example-token",
	}
	for _, item := range api.Headers {
		headers[item.Name] = item.Value
	}
	err := callAPI(api.Method, api.URL, headers, api.Params, data)
	if err != nil {
		return err
	}
	if mode == models.AppModeInput {
	} else if mode == models.AppModeOutPut {
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
