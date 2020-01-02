package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Configuration struct {
	InstanceURL *string

	InstanceToken *string
	AaiToken      *string

	ChunkSize *int64
}

func createConfiguration(instanceURL string) {
	configuration := Configuration{}
	configuration.InstanceURL = &instanceURL
	var defaultChunkSize int64 = 200
	configuration.ChunkSize = &defaultChunkSize
	saveConfiguration(configuration)
}

func loadConfiguration() (*Configuration, error) {
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

func saveConfiguration(configuration Configuration) {
	userCacheDir, _ := os.UserCacheDir()
	bytes, err := json.Marshal(configuration)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join(userCacheDir, "lega-uploader", "config.json"), bytes, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
