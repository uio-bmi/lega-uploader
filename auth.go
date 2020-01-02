package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
	"syscall"
)

func login() {
	configuration, _ := loadConfiguration()
	println("Username: ")
	var username string
	_, _ = fmt.Scanln(&username)
	println("Password: ")
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	request, _ := http.NewRequest("GET", *configuration.InstanceURL+"/cega", nil)
	request.SetBasicAuth(username, string(bytePassword))
	response, err := http.DefaultClient.Do(request)
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
