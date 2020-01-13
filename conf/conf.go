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
	_, _ = fmt.Scanln(&instanceURL)
	configurationProvider := NewConfigurationProvider()
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
}

func NewConfigurationProvider() ConfigurationProvider {
	return defaultConfigurationProvider{}
}

func (cp defaultConfigurationProvider) LoadConfiguration() (*Configuration, error) {
	userCacheDir, _ := os.UserCacheDir()
	configFile, err := os.Open(filepath.Join(userCacheDir, "lega-uploader", "config.json"))
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
	userCacheDir, _ := os.UserCacheDir()
	bytes, err := json.Marshal(configuration)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(userCacheDir, "lega-uploader", "config.json"), bytes, os.ModePerm); err != nil {
		return err
	}
	return nil
}
