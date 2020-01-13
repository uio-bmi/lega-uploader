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

type Configuration struct {
	InstanceURL *string

	InstanceToken *string
	AaiToken      *string

	ChunkSize *int64
}

func Configure() {
	fmt.Println(aurora.Yellow("Instance URL: "))
	var instanceURL string
	_, _ = fmt.Scanln(&instanceURL)
	if err := createConfiguration(strings.TrimRight(instanceURL, "/")); err != nil {
		log.Fatal(aurora.Red(err))
	}
}

func createConfiguration(instanceURL string) error {
	configuration := Configuration{}
	configuration.InstanceURL = &instanceURL
	var defaultChunkSize int64 = 200
	configuration.ChunkSize = &defaultChunkSize
	if err := SaveConfiguration(configuration); err != nil {
		return err
	}
	return nil
}

func LoadConfiguration() (*Configuration, error) {
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

func SaveConfiguration(configuration Configuration) error {
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
