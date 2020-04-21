// Package resuming contains structures and methods for listing or resuming uploads.
package resuming

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/requests"
	"io/ioutil"
	"net/http"
	"strings"
)

// Resumable structure represents resumable upload.
type Resumable struct {
	ID    string
	Name  string
	Size  int64
	Chunk int64
}

// ResumablesManager interface provides method for listing resumable uploads.
type ResumablesManager interface {
	GetResumables() (*[]Resumable, error)
}

type defaultResumablesManager struct {
	client requests.Client
}

// NewResumablesManager constructs ResumablesManager using requests.Client.
func NewResumablesManager(client *requests.Client) (ResumablesManager, error) {
	resumablesManager := defaultResumablesManager{}
	if client != nil {
		resumablesManager.client = *client
	} else {
		resumablesManager.client = requests.NewClient(nil)
	}
	return resumablesManager, nil
}

// GetResumables method lists resumable uploads.
func (rm defaultResumablesManager) GetResumables() (*[]Resumable, error) {
	configuration := conf.NewConfiguration()
	response, err := rm.client.DoRequest(http.MethodGet,
		configuration.GetLocalEGAInstanceURL()+"/resumables",
		nil,
		map[string]string{"Proxy-Authorization": "Bearer " + configuration.GetElixirAAIToken()},
		nil,
		configuration.GetCentralEGAUsername(),
		configuration.GetCentralEGAPassword())
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = response.Body.Close()
	if err != nil {
		return nil, err
	}
	resumables := make([]Resumable, 0)
	_, err = jsonparser.ArrayEach(body,
		func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			id, _ := jsonparser.GetString(value, "id")
			fileName, _ := jsonparser.GetString(value, "fileName")
			fileName = fileName[strings.Index(fileName, "-")+1:]
			size, _ := jsonparser.GetInt(value, "nextOffset")
			chunk, _ := jsonparser.GetInt(value, "maxChunk")
			chunk++
			resumable := Resumable{id, fileName, size, chunk}
			resumables = append(resumables, resumable)
		},
		"resumables")
	if err != nil {
		return nil, err
	}
	return &resumables, nil
}
