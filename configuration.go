package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	InstanceURL *string

	InstanceToken *string
	AaiToken      *string

	ChunkSize *int64
}

func configure() {
	println("Instance URL: ")
	var instanceURL string
	_, _ = fmt.Scanln(&instanceURL)
	createConfiguration(strings.TrimRight(instanceURL, "/"))
}

func createConfiguration(instanceURL string) {
	configuration := Configuration{}
	configuration.InstanceURL = &instanceURL
	var defaultChunkSize int64 = 200
	configuration.ChunkSize = &defaultChunkSize
	saveConfiguration(configuration)
}

func loadConfiguration() *Configuration {
	userCacheDir, _ := os.UserCacheDir()
	configFile, err := os.Open(filepath.Join(userCacheDir, "lega-uploader", "config.json"))
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()
	configFileContent, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var configuration Configuration
	if err = json.Unmarshal(configFileContent, &configuration); err != nil {
		log.Fatal(err)
	}
	return &configuration
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
