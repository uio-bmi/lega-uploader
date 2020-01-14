package conf

import (
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var configFile *os.File

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	var err error
	configFile, err = ioutil.TempFile(".", "config.json")
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}

func TestNewConfigurationInstanceURL(t *testing.T) {
	configuration := NewConfiguration("test", nil)
	if *configuration.InstanceURL != "test" {
		t.Error()
	}
}

func TestNewConfigurationEmptyChunkSize(t *testing.T) {
	configuration := NewConfiguration("test", nil)
	if *configuration.ChunkSize != 200 {
		t.Error()
	}
}

func TestNewConfigurationFilledChunkSize(t *testing.T) {
	defaultChunkSize := int64(1)
	configuration := NewConfiguration("test", &defaultChunkSize)
	if *configuration.ChunkSize != defaultChunkSize {
		t.Error()
	}
}

func TestLoadConfigurationAbsent(t *testing.T) {
	filePath := "test"
	configurationProvider, err := NewConfigurationProvider(&filePath)
	if err != nil {
		t.Error(err)
	}
	_, err = configurationProvider.LoadConfiguration()
	if err == nil {
		t.Error(err)
	}
}

func TestLoadConfigurationMalformed(t *testing.T) {
	filePath := "conf.go"
	configurationProvider, err := NewConfigurationProvider(&filePath)
	if err != nil {
		t.Error(err)
	}
	_, err = configurationProvider.LoadConfiguration()
	if err == nil {
		t.Error(err)
	}
}

func TestLoadConfigurationSuccess(t *testing.T) {
	filePath := "conf.go"
	configurationProvider, err := NewConfigurationProvider(&filePath)
	if err != nil {
		t.Error(err)
	}
	_, err = configurationProvider.LoadConfiguration()
	if err == nil {
		t.Error(err)
	}
}

func teardown() {
	err := os.Remove(configFile.Name())
	if err != nil {
		log.Fatal(aurora.Red(err))
	}
}
