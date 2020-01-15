package conf

import (
	"encoding/json"
	"fmt"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Configure() {
	fmt.Println(aurora.Yellow("Instance URL: "))
	var instanceURL string
	_, err := fmt.Scanln(&instanceURL)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	configurationProvider, err := NewConfigurationProvider(nil)
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
	configuration := NewConfiguration(strings.TrimRight(instanceURL, "/"), nil)
	if err := configurationProvider.SaveConfiguration(configuration); err != nil {
		log.Fatal(aurora.Red(err))
	}
}

type Configuration struct {
	InstanceURL *string

	InstanceToken *string
	AaiToken      *string

	ChunkSize *int64
}

func NewConfiguration(instanceURL string, chunkSize *int64) Configuration {
	configuration := Configuration{InstanceURL: &instanceURL}
	if chunkSize != nil {
		configuration.ChunkSize = chunkSize
	} else {
		var defaultChunkSize int64 = 200
		configuration.ChunkSize = &defaultChunkSize
	}
	return configuration
}

type ConfigurationProvider interface {
	LoadConfiguration() (*Configuration, error)
	SaveConfiguration(configuration Configuration) error
}

type defaultConfigurationProvider struct {
	configFile string
}

func NewConfigurationProvider(configFile *string) (ConfigurationProvider, error) {
	configurationProvider := defaultConfigurationProvider{}
	if configFile != nil {
		configurationProvider.configFile = *configFile
	} else {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		configurationProvider.configFile = filepath.Join(userCacheDir, "lega-uploader", "config.json")
	}
	return configurationProvider, nil
}

func (cp defaultConfigurationProvider) LoadConfiguration() (*Configuration, error) {
	configFile, err := os.Open(cp.configFile)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	configFileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	var configuration Configuration
	if err = json.Unmarshal(configFileContent, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func (cp defaultConfigurationProvider) SaveConfiguration(configuration Configuration) error {
	bytes, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(cp.configFile, bytes, os.ModePerm); err != nil {
		return err
	}
	return nil
}
