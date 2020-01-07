package main

import (
	"errors"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Resumable struct {
	id    string
	name  string
	size  int64
	chunk int64
}

func resumables() {
	resumables, err := getResumables()
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	for _, resumable := range *resumables {
		fmt.Println(aurora.Blue(resumable.name + "\t (" + strconv.FormatInt(resumable.size, 10) + " bytes uploaded)"))
	}
}

func getResumables() (*[]Resumable, error) {
	configuration := loadConfiguration()
	response, err := doRequest(http.MethodGet,
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
