package db

import (
	"bytes"
	"encoding/json"
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
