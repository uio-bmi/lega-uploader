package main

import (
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//func resume() error {
//	return nil
//}

func resumables() {
	configuration, _ := loadConfiguration()
	response, err := doRequest(http.MethodGet,
		*configuration.InstanceURL+"/resumables",
		nil,
		map[string]string{"Authorization": "Bearer " + *configuration.InstanceToken},
		nil,
		nil,
		nil)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal(response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	_, err = jsonparser.ArrayEach(body,
		func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			fileName, _ := jsonparser.GetString(value, "fileName")
			size, _ := jsonparser.GetInt(value, "nextOffset")
			println(fileName[strings.Index(fileName, "-")+1:] + "\t (" + strconv.FormatInt(size, 10) + " bytes uploaded)")
		},
		"resumables")
	if err != nil {
		log.Fatal(err)
	}
}
