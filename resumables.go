package main

import (
	"errors"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Resumable struct {
	name string
	size int64
}

func resume(path string) error {
	println(path)
	resumables, err := getResumbales()
	if err != nil {
		return err
	}
	for _, resumable := range *resumables {
		println(resumable.name)
	}
	return nil
}

func resumables() {
	resumables, err := getResumbales()
	if err != nil {
		log.Fatal(err)
	}
	for _, resumable := range *resumables {
		println(resumable.name[strings.Index(resumable.name, "-")+1:] + "\t (" + strconv.FormatInt(resumable.size, 10) + " bytes uploaded)")
	}
}

func getResumbales() (*[]Resumable, error) {
	configuration, _ := loadConfiguration()
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
			fileName, _ := jsonparser.GetString(value, "fileName")
			size, _ := jsonparser.GetInt(value, "nextOffset")
			resumable := Resumable{fileName, size}
			resumables = append(resumables, resumable)
		},
		"resumables")
	if err != nil {
		return nil, err
	}
	return &resumables, nil
}
