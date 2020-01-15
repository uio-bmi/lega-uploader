package requests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var server *httptest.Server
var client Client
var url string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, password, ok := r.BasicAuth()
		println(auth)
		println(password)
		println(ok)
		w.WriteHeader(200)
	}))
	client = NewClient(server.Client())
	url = server.URL
}

func TestDoRequest(t *testing.T) {
	get, err := client.DoRequest(http.MethodGet, url, nil, nil, nil, nil, nil)
	if err != nil {
		t.Error()
	}
	println(get)
}

func teardown() {
	server.Close()
}
