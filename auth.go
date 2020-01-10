package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"net/http"
)

func login(username string, password string) error {
	configuration := loadConfiguration()
	response, err := doRequest(http.MethodGet, *configuration.InstanceURL+"/cega", nil, nil, nil, &username, &password)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return errors.New("Authentication failure")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	token := string(body)
	if "Wrong credentials" == token {
		return errors.New("Wrong credentials")
	}
	configuration.InstanceToken = &token
	saveConfiguration(*configuration)
	fmt.Println(aurora.Green("Success!"))
	return nil
}
