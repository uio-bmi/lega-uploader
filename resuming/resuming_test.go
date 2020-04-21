package resuming

import (
	"github.com/uio-bmi/lega-uploader/requests"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

type mockClient struct {
}

func (mockClient) DoRequest(_ string, url string, _ io.Reader, _ map[string]string, _ map[string]string, _ string, _ string) (*http.Response, error) {
	if strings.HasSuffix(url, "/resumables") {
		body := ioutil.NopCloser(strings.NewReader(`{"resumables": [{"id": "1", "fileName": "dummy-test.enc", "nextOffset": 100, "maxChunk": 10}]}`))
		response := http.Response{StatusCode: 200, Body: body}
		return &response, nil
	}
	return nil, nil
}

func TestGetResumables(t *testing.T) {
	_ = os.Setenv("CENTRAL_EGA_USERNAME", "user")
	_ = os.Setenv("CENTRAL_EGA_PASSWORD", "pass")
	_ = os.Setenv("LOCAL_EGA_INSTANCE_URL", "http://localhost/")
	_ = os.Setenv("ELIXIR_AAI_TOKEN", "token")
	var client requests.Client = mockClient{}
	resumablesManager, err := NewResumablesManager(&client)
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
