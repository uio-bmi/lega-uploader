package auth

import (
	"../conf"
	"../requests"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockConfigurationProvider struct {
}

func (cp mockConfigurationProvider) LoadConfiguration() (*conf.Configuration, error) {
	configuration := conf.NewConfiguration("http://localhost/", nil)
	return &configuration, nil
}

func (cp mockConfigurationProvider) SaveConfiguration(_ conf.Configuration) error {
	return nil
}

type mockClient struct {
}

func (mockClient) DoRequest(_ string, url string, _ io.Reader, _ map[string]string, _ map[string]string, username *string, password *string) (*http.Response, error) {
	if strings.HasSuffix(url, "/cega") {
		if username == nil || password == nil {
			body := ioutil.NopCloser(strings.NewReader("Wrong credentials"))
			response := http.Response{StatusCode: 403, Body: body}
			return &response, nil
		}
		if *username == "dummy" && *password == "dummy" {
			body := ioutil.NopCloser(strings.NewReader("Success"))
			response := http.Response{StatusCode: 200, Body: body}
			return &response, nil
		}
		body := ioutil.NopCloser(strings.NewReader("Wrong credentials"))
		response := http.Response{StatusCode: 403, Body: body}
		return &response, nil
	}
	return nil, nil
}

func TestLoginSuccess(t *testing.T) {
	var configurationProvider conf.ConfigurationProvider = mockConfigurationProvider{}
	var client requests.Client = mockClient{}
	authenticationManager := NewAuthenticationManager(&configurationProvider, &client)
	err := authenticationManager.Authenticate("dummy", "dummy")
	if err != nil {
		t.Error(err)
	}
}

func TestLoginFailure(t *testing.T) {
	var client requests.Client = mockClient{}
	authenticationManager := NewAuthenticationManager(nil, &client)
	err := authenticationManager.Authenticate("", "")
	if err == nil {
		t.Error()
	}
}
