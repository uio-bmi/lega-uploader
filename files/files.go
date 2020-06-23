// Package files contains structures and methods for listing or deleting uploaded files.
package files

import (
	"errors"
	"github.com/buger/jsonparser"
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/requests"
	"io/ioutil"
	"net/http"
)

// File structure represents uploaded File.
type File struct {
	FileName     string
	Size         int64
	ModifiedDate string
}

// FileManager interface provides method for managing uploaded files.
type FileManager interface {
	ListFiles() (*[]File, error)
	DeleteFile(fileName string) error
}

type defaultFileManager struct {
	client requests.Client
}

// NewFileManager constructs FileManager using requests.Client.
func NewFileManager(client *requests.Client) (FileManager, error) {
	fileManager := defaultFileManager{}
	if client != nil {
		fileManager.client = *client
	} else {
		fileManager.client = requests.NewClient(nil)
	}
	return fileManager, nil
}

// ListFiles method lists uploaded files.
func (rm defaultFileManager) ListFiles() (*[]File, error) {
	configuration := conf.NewConfiguration()
	response, err := rm.client.DoRequest(http.MethodGet,
		configuration.GetLocalEGAInstanceURL()+"/files",
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
	files := make([]File, 0)
	_, _ = jsonparser.ArrayEach(body,
		func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			fileName, _ := jsonparser.GetString(value, "fileName")
			size, _ := jsonparser.GetInt(value, "size")
			modifiedDate, _ := jsonparser.GetString(value, "modifiedDate")
			file := File{fileName, size, modifiedDate}
			files = append(files, file)
		},
		"files")
	return &files, nil
}

// DeleteFiles method deletes uploaded file by its name.
func (rm defaultFileManager) DeleteFile(fileName string) error {
	configuration := conf.NewConfiguration()
	response, err := rm.client.DoRequest(http.MethodDelete,
		configuration.GetLocalEGAInstanceURL()+"/files",
		nil,
		map[string]string{"Proxy-Authorization": "Bearer " + configuration.GetElixirAAIToken()},
		map[string]string{"fileName": fileName},
		configuration.GetCentralEGAUsername(),
		configuration.GetCentralEGAPassword())
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	return nil
}
