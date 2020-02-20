// Package conf contains structure and methods to create/load configuration for the application.
package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Configuration structure is a holder for application settings.
type Configuration struct {
	InstanceURL *string

	InstanceToken *string
	AaiToken      *string

	ChunkSize *int64
}

// NewConfiguration constructs Configuration, accepting LocalEGA URL instance and possibly chunk size.
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

// ConfigurationProvider is an interface providing methods for reading and writing configuration.
type ConfigurationProvider interface {
	LoadConfiguration() (*Configuration, error)
	SaveConfiguration(configuration Configuration) error
}

type defaultConfigurationProvider struct {
	configFile string
}

// NewConfigurationProvider constructs ConfigurationProvider from a path to the config-file.
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

// LoadConfiguration methods loads the configuration from a file, returning Configuration structure.
func (cp defaultConfigurationProvider) LoadConfiguration() (*Configuration, error) {
	configFile, err := os.Open(cp.configFile)
	if err != nil {
		return nil, err
	}
	configFileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	err = configFile.Close()
	if err != nil {
		return nil, err
	}
	var configuration Configuration
	if err = json.Unmarshal(configFileContent, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

// Save configuration serializes Configuration structure to a file.
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
