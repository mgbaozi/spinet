package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Query(method string, url string, headers map[string]string, params interface{}, response interface{}) error {
	var req *http.Request
	var err error
	if headers == nil {
		headers = make(map[string]string)
	}
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
