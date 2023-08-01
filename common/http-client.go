package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Response[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func HttpGet[T any](url string, token ...interface{}) (*Response[T], error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "请求参数异常",
		}, err
	}

	// 添加Cookie到请求头
	if len(token) > 0 {
		cookie := &http.Cookie{
			Name:  "pushToken",
			Value: token[0].(string),
		}
		req.AddCookie(cookie)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "请求失败，请稍后重试！",
		}, err
	}
	defer res.Body.Close()
	var response Response[T]
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "响应内容解析异常",
		}, err
	}

	return &response, nil
}

func HttpPost[T any](url string, postData interface{}) (*Response[T], error) {
	buffer := bytes.NewBuffer(nil)
	buffer = nil

	if postData != nil {
		jsonData, err := json.Marshal(postData)
		if err != nil {
			return &Response[T]{
				Code: 500,
				Msg:  "请求参数异常",
			}, err
		}
		buffer = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "请求参数异常",
		}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "请求失败，请稍后重试！",
		}, err
	}
	defer res.Body.Close()
	var response Response[T]
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return &Response[T]{
			Code: 500,
			Msg:  "响应内容解析异常",
		}, err
	}

	return &response, nil
}
