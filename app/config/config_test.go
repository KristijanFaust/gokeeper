package config

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// LoadConfiguration should successfully load to ApplicationConfig values from the default .yml file,
// and the current state of the configuration file should be valid.
func TestLoadConfiguration(t *testing.T) {
	config := LoadConfiguration("../../config.yml")
	assert.NotNil(t, config, "Application configuration should be setup now")
}

// LoadConfiguration should panic on any error while reading from the given configuration file
func TestLoadConfigurationWithConfigurationReadError(t *testing.T) {
	assert.PanicsWithValue(
		t, "Error occured while trying to read configuration file: open nonexistent-configuration-file: no such file or directory",
		func() { LoadConfiguration("nonexistent-configuration-file") },
		"LoadConfiguration should panic when passed a nonexistent configuration file path",
	)
}

// LoadConfiguration should panic on any error while parsing values from the given configuration file
func TestLoadConfigurationWithInvalidConfigurationError(t *testing.T) {
	generateInvalidConfiguration()
	defer removeInvalidConfiguration()
	assert.PanicsWithValue(
		t, "Error occured while trying to decode configuration values: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid...` into config.Config",
		func() { LoadConfiguration("./invalid-config.yml") },
		"LoadConfiguration should panic when passed a nonexistent configuration file path",
	)
}

func generateInvalidConfiguration() {
	configurationData := []byte("invalid configuration")
	err := ioutil.WriteFile("./invalid-config.yml", configurationData, 0644)
	if err != nil {
		log.Panic("Could not generate invalid configuration file")
	}
}

func removeInvalidConfiguration() {
	err := os.Remove("./invalid-config.yml")
	if err != nil {
		log.Println("Could not delete invalid-config.yml. Please consider removing it manually if needed.")
		log.Println("Error: ", err)
	}
}
