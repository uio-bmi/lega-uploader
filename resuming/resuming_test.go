package resuming

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

func (mockClient) DoRequest(_ string, url string, _ io.Reader, _ map[string]string, _ map[string]string, _ *string, _ *string) (*http.Response, error) {
	if strings.HasSuffix(url, "/resumables") {
		body := ioutil.NopCloser(strings.NewReader(`{"resumables": [{"id": "1", "fileName": "dummy-test.enc", "nextOffset": 100, "maxChunk": 10}]}`))
		response := http.Response{StatusCode: 200, Body: body}
		return &response, nil
	}
	return nil, nil
}

func TestGetResumables(t *testing.T) {
	var configurationProvider conf.ConfigurationProvider = &mockConfigurationProvider{}
	newConfiguration := conf.NewConfiguration("http://localhost/", nil)
	token := "token"
	newConfiguration.InstanceToken = &token
	err := configurationProvider.SaveConfiguration(newConfiguration)
	if err != nil {
		t.Error(err)
	}
	var client requests.Client = mockClient{}
	resumablesManager, err := NewResumablesManager(&configurationProvider, &client)
	if err != nil {
		t.Error(err)
	}
	resumables, err := resumablesManager.GetResumables()
	if err != nil {
		t.Error(err)
	}
	if resumables == nil || len(*resumables) != 1 {
		t.Error()
	}
	resumable := (*resumables)[0]
	if resumable.ID != "1" || resumable.Name != "test.enc" || resumable.Size != 100 || resumable.Chunk != 11 {
		t.Error()
	}
}
