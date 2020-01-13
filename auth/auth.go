package auth

import (
	"../conf"
	"../requests"
	"errors"
	"fmt"
	"github.com/logrusorgru/aurora"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net/http"
)

func Login() {
	fmt.Println(aurora.Yellow("Username: "))
	var username string
	_, _ = fmt.Scanln(&username)
	fmt.Println(aurora.Yellow("Password: "))
	bytePassword, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	password := string(bytePassword)
	authenticationManager := NewAuthenticationManager(nil)
	if err := authenticationManager.Authenticate(username, password); err != nil {
		log.Fatal(aurora.Red(err))
	} else {
		fmt.Println(aurora.Green("Success!"))
	}
}

type AuthenticationManager interface {
	Authenticate(username string, password string) error
}

type defaultAuthenticationManager struct {
	Client requests.Client
}

func NewAuthenticationManager(client *requests.Client) AuthenticationManager {
	if client != nil {
		return defaultAuthenticationManager{Client: *client}
	} else {
		return defaultAuthenticationManager{Client: requests.NewClient()}
	}
}

func (am defaultAuthenticationManager) Authenticate(username string, password string) error {
	configuration, err := conf.LoadConfiguration()
	if err != nil {
		return err
	}
	response, err := am.Client.DoRequest(http.MethodGet, *configuration.InstanceURL+"/cega", nil, nil, nil, &username, &password)
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
	if err = conf.SaveConfiguration(*configuration); err != nil {
		return err
	}
	return nil
}
