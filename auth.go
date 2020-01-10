package main

import (
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
)

func login() error {
	configuration := loadConfiguration()
	fmt.Println(aurora.Yellow("Username: "))
	var username string
	_, _ = fmt.Scanln(&username)
	fmt.Println(aurora.Yellow("Password: "))
	bytePassword, err := terminal.ReadPassword(0)
	password := string(bytePassword)
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
