package main

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
)

func login() {
	configuration := loadConfiguration()
	fmt.Println(aurora.Yellow("Username: "))
	var username string
	_, _ = fmt.Scanln(&username)
	fmt.Println(aurora.Yellow("Password: "))
	bytePassword, err := terminal.ReadPassword(0)
	password := string(bytePassword)
	response, err := doRequest(http.MethodGet, *configuration.InstanceURL+"/cega", nil, nil, nil, &username, &password)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatal(aurora.Red("Authentication failure"))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	token := string(body)
	if "Wrong credentials" == token {
		log.Fatal(aurora.Red("Wrong credentials"))
	}
	configuration.InstanceToken = &token
	saveConfiguration(*configuration)
	fmt.Println(aurora.Green("Success!"))
}
