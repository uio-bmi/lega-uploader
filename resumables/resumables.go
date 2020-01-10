package resumables

import (
	"../conf"
	"../requests"
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Resumable struct {
	Id    string
	Name  string
	Size  int64
	Chunk int64
}

func Resumables() error {
	resumables, err := GetResumables()
	if err != nil {
		return err
	}
	for _, resumable := range *resumables {
		fmt.Println(aurora.Blue(resumable.Name + "\t (" + strconv.FormatInt(resumable.Size, 10) + " bytes uploaded)"))
	}
	return nil
}

func GetResumables() (*[]Resumable, error) {
	configuration := conf.LoadConfiguration()
	response, err := requests.DoRequest(http.MethodGet,
		*configuration.InstanceURL+"/resumables",
		nil,
		map[string]string{"Authorization": "Bearer " + *configuration.InstanceToken},
		nil,
		nil,
		nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
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
