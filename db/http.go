package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Http[T any] struct {
	url string
}

func NewHttp[T any](url string) *Http[T] {

	return &Http[T]{url: url}
}

func (h *Http[T]) GetText(param url.Values) (string, error) {

	fullURL := h.url + "?" + param.Encode()

	r, err := http.Get(fullURL)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil

}

func (h *Http[T]) Get(param url.Values) (*T, error) {

	fullURL := h.url + "?" + param.Encode()

	r, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	var data T

	err = json.Unmarshal(bodyBytes, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

func (h *Http[T]) Post(param url.Values, body map[string]string) (*T, error) {

	fullURL := h.url + "?" + param.Encode()

	d, err := json.Marshal(body)

	if err != nil {
		return nil, err
	}

	r, err := http.Post(fullURL, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	var data T

	err = json.Unmarshal(bodyBytes, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}

// PasteUrl 拼接url
func PasteUrl(baseURL string, params map[string]string) string {
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Add(key, value)
	}
	return baseURL + "?" + queryParams.Encode()
}

// GetUrl 获取url的数据
func GetUrl[T any](url string) (T, error) {
	var t T
	resp, err := http.Get(url)
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}

// PostUrl 发送post请求
func PostUrl[T any](url string, data map[string]any) (T, error) {
	var t T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return t, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return t, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return t, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}
