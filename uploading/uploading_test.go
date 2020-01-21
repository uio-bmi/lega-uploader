package uploading

import (
	"../conf"
	"../requests"
	"../resuming"
	"fmt"
	"github.com/logrusorgru/aurora"
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
	teardown()
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
	dir, err = ioutil.TempDir(".", "test")
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	file, err = ioutil.TempFile(dir, "test.dat")
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	err = file.Truncate(1e7)
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
			if checksum != "5f363e0e58a95f06cbe9bbc662c5dfb6" {
				body := ioutil.NopCloser(strings.NewReader(""))
				response = http.Response{StatusCode: 500, Body: body}
				return &response, nil
			}
		} else if chunk == "2" {
			checksum := params["md5"]
			if checksum != "443b86c26729f5f24407c7236957d653" {
				body := ioutil.NopCloser(strings.NewReader(""))
				response = http.Response{StatusCode: 500, Body: body}
				return &response, nil
			}
		} else if chunk == "end" {
			checksum := params["md5"]
			fileSize := params["fileSize"]
			if checksum != "d41d8cd98f00b204e9800998ecf8427e" || fileSize != "10000000" {
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
	if err != nil {
		t.Error(err)
	}
}

func teardown() {
	err := os.Remove(file.Name())
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	err = os.Remove(dir)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}
