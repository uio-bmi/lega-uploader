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
	authenticationManager, err := NewAuthenticationManager(nil, nil)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
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
	configurationProvider conf.ConfigurationProvider
	client                requests.Client
}

func NewAuthenticationManager(configurationProvider *conf.ConfigurationProvider, client *requests.Client) (AuthenticationManager, error) {
	authenticationManager := defaultAuthenticationManager{}
	if configurationProvider != nil {
		authenticationManager.configurationProvider = *configurationProvider
	} else {
		newConfigurationProvider, err := conf.NewConfigurationProvider(nil)
		if err != nil {
			return nil, err
		}
		authenticationManager.configurationProvider = newConfigurationProvider
	}
	if client != nil {
		authenticationManager.client = *client
	} else {
		authenticationManager.client = requests.NewClient()
	}
	return authenticationManager, nil
}

func (am defaultAuthenticationManager) Authenticate(username string, password string) error {
	configuration, err := am.configurationProvider.LoadConfiguration()
	if err != nil {
		return err
	}
	response, err := am.client.DoRequest(http.MethodGet, *configuration.InstanceURL+"/cega", nil, nil, nil, &username, &password)
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
	if err = am.configurationProvider.SaveConfiguration(*configuration); err != nil {
		return err
	}
	return nil
}
