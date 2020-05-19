package uploading

import (
	"fmt"
	"github.com/logrusorgru/aurora"
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
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	_ = os.Setenv("CENTRAL_EGA_USERNAME", "user")
	_ = os.Setenv("CENTRAL_EGA_PASSWORD", "pass")
	_ = os.Setenv("LOCAL_EGA_INSTANCE_URL", "http://localhost/")
	_ = os.Setenv("ELIXIR_AAI_TOKEN", "token")
	var client requests.Client = mockClient{}
	resumablesManager, err := resuming.NewResumablesManager(&client)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	uploader, err = NewUploader(&client, &resumablesManager)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	dir = "../test"
	file, err = os.Open("../test/sample.txt.enc")
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}

type mockClient struct {
}

func (mockClient) DoRequest(_, url string, _ io.Reader, headers, params map[string]string, username, password string) (*http.Response, error) {
	var response http.Response
	if username != "user" || password != "pass" || !strings.HasPrefix(headers["Proxy-Authorization"], "Bearer ") {
		body := ioutil.NopCloser(strings.NewReader(""))
		response = http.Response{StatusCode: 401, Body: body}
		return &response, nil
	}
	if strings.Contains(url, "/stream") {
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
	if err == nil || !strings.HasSuffix(err.Error(), "not a Crypt4GH file") {
		t.Error(err)
	}
}

func teardown() {
	_ = os.Unsetenv("CENTRAL_EGA_USERNAME")
	_ = os.Unsetenv("CENTRAL_EGA_PASSWORD")
	_ = os.Unsetenv("LOCAL_EGA_INSTANCE_URL")
	_ = os.Unsetenv("ELIXIR_AAI_TOKEN")
}
