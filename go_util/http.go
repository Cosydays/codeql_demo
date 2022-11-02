package go_util

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func HttpGet(ctx context.Context, reqUrl string, params map[string]string, header map[string]string) (string, error) {
	return HttpGetOfAuth(ctx, reqUrl, params, header, nil)
}

func HttpGetOfAuth(ctx context.Context, reqUrl string, params map[string]string, header map[string]string, auth map[string]string) (string, error) {
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	req.URL.RawQuery = query.Encode()
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	if auth != nil {
		for k, v := range auth {
			req.SetBasicAuth(k, v)
		}
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = httpClient.Do(req)
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+2))
	}

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("status code: %v, message: %v", resp.StatusCode, string(body)))
	}
	return string(body), nil
}

func HttpPost(ctx context.Context, reqUrl string, jsonParams []byte, header map[string]string, timeout time.Duration) (string, error) {
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonParams))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	httpClient := &http.Client{Timeout: timeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("status code: %v, message: %v", resp.StatusCode, string(body)))
	}
	return string(body), nil
}
