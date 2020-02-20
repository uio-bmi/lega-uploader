package auth

import (
	"errors"
	"github.com/uio-bmi/lega-uploader/conf"
	"github.com/uio-bmi/lega-uploader/requests"
	"io/ioutil"
	"net/http"
)

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
		authenticationManager.client = requests.NewClient(nil)
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
	if response.StatusCode != 200 {
		return errors.New("Authentication failure")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = response.Body.Close()
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
