package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
)

func login() {
	configuration, _ := loadConfiguration()
	println("Username: ")
	var username string
	_, _ = fmt.Scanln(&username)
	println("Password: ")
	bytePassword, err := terminal.ReadPassword(0)
	password := string(bytePassword)
	response, err := doRequest(http.MethodGet, *configuration.InstanceURL+"/cega", nil, nil, nil, &username, &password)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal("Authentication failure")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	token := string(body)
	if "Wrong credentials" == token {
		log.Fatal("Wrong credentials")
	}
	configuration.InstanceToken = &token
	saveConfiguration(*configuration)
	println("Success!")
}
