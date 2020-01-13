package requests

import (
	"io"
	"net/http"
)

type Client interface {
	DoRequest(method string, url string, body io.Reader, headers map[string]string, params map[string]string, username *string, password *string) (*http.Response, error)
}

type defaultClient struct {
}

func NewClient() Client {
	return defaultClient{}
}

func (defaultClient) DoRequest(method string, url string, body io.Reader, headers map[string]string, params map[string]string, username *string, password *string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for name, header := range headers {
			request.Header.Add(name, header)
		}
	}
	if params != nil {
		query := request.URL.Query()
		for name, param := range params {
			query.Add(name, param)
		}
		request.URL.RawQuery = query.Encode()
	}
	if username != nil && password != nil {
		request.SetBasicAuth(*username, *password)
	}
	return http.DefaultClient.Do(request)
}
