package uploading

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/requests"
	"github.com/uio-bmi/lega-uploader/resuming"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var uploader Uploader
var dir string
var file *os.File

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	var configurationProvider conf.ConfigurationProvider = &mockConfigurationProvider{}
	var defaultChunk = int64(5)
	newConfiguration := conf.NewConfiguration("http://localhost/", &defaultChunk)
	token := "token"
	newConfiguration.InstanceToken = &token
	err = configurationProvider.SaveConfiguration(newConfiguration)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	var client requests.Client = mockClient{}
	resumablesManager, err := resuming.NewResumablesManager(&configurationProvider, &client)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	uploader, err = NewUploader(&configurationProvider, &client, &resumablesManager)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	dir = "../test"
	file, err = os.Open("../test/sample.txt.enc")
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}

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

func (mockClient) DoRequest(_ string, url string, _ io.Reader, _ map[string]string, params map[string]string, _ *string, _ *string) (*http.Response, error) {
	if strings.Contains(url, "/stream") {
		var response http.Response
		fmt.Printf("%v\n", params)
		var id string
		var ok bool
		id, ok = params["id"]
		if !ok {
			id = "123"
		}
		if id != "123" {
			body := ioutil.NopCloser(strings.NewReader(""))
			response = http.Response{StatusCode: 500, Body: body}
			return &response, nil
		}
		chunk := params["chunk"]
		if chunk == "1" {
			checksum := params["md5"]
			if checksum != "da385d93ae510bc91c9c8af7e670ac6f" {
				body := ioutil.NopCloser(strings.NewReader(""))
				response = http.Response{StatusCode: 500, Body: body}
				return &response, nil
			}
		} else if chunk == "end" {
			checksum := params["md5"]
			fileSize := params["fileSize"]
			if checksum != "d41d8cd98f00b204e9800998ecf8427e" || fileSize != "65688" {
				body := ioutil.NopCloser(strings.NewReader(""))
				response = http.Response{StatusCode: 500, Body: body}
				return &response, nil
			}
		}
		body := ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`{"id":"%s"}`, id)))
		response = http.Response{StatusCode: 200, Body: body}
		return &response, nil
	}
	return nil, nil
}

func TestUploadFile(t *testing.T) {
	err := uploader.Upload(file.Name(), false)
	if err != nil {
		t.Error(err)
	}
}

func TestUploadFolder(t *testing.T) {
	err := uploader.Upload(dir, false)
	if err == nil || err.Error() != "not a Crypt4GH file" {
		t.Error(err)
	}
}
