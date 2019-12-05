package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
	"syscall"
)

func main() {
	initialize()

	loginFlag := flag.Bool("login", false, "Logs you in.")
	flag.Parse()
	if *loginFlag {
		login()
	}
}

func initialize() {
	if _, err := loadConfiguration(); err != nil {
		createConfiguration("https://localhost")
	}
}

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
	all, err := ioutil.ReadAll(response.Body)
	println(string(all))
	println(err)
}
