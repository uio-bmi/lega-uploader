package auth

import (
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/requests"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type mockConfigurationProvider struct {
	configuration *conf.Configuration
}

func (cp mockConfigurationProvider) LoadConfiguration() (*conf.Configuration, error) {
	return cp.configuration, nil
}

func (cp *mockConfigurationProvider) SaveConfiguration(configuration conf.Configuration) error {
	cp.configuration = &configuration
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

func TestAuthenticateSuccess(t *testing.T) {
	var configurationProvider conf.ConfigurationProvider = &mockConfigurationProvider{}
	newConfiguration := conf.NewConfiguration("http://localhost/", nil)
	err := configurationProvider.SaveConfiguration(newConfiguration)
	if err != nil {
		t.Error(err)
	}
	var client requests.Client = mockClient{}
	authenticationManager, err := NewAuthenticationManager(&configurationProvider, &client)
	if err != nil {
		t.Error(err)
	}
	err = authenticationManager.Authenticate("dummy", "dummy")
	if err != nil {
		t.Error(err)
	}
	configuration, err := configurationProvider.LoadConfiguration()
	if err != nil {
		t.Error(err)
	}
	if *configuration.InstanceToken != "Success" {
		t.Error()
	}
}

func TestAuthenticateFailure(t *testing.T) {
	var configurationProvider conf.ConfigurationProvider = &mockConfigurationProvider{}
	newConfiguration := conf.NewConfiguration("http://localhost/", nil)
	err := configurationProvider.SaveConfiguration(newConfiguration)
	if err != nil {
		t.Error(err)
	}
	var client requests.Client = mockClient{}
	authenticationManager, err := NewAuthenticationManager(&configurationProvider, &client)
	if err != nil {
		t.Error(err)
	}
	err = authenticationManager.Authenticate("", "")
	if err == nil {
		t.Error()
	}
}
