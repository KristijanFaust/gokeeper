package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var ApplicationConfig *Config

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Datasource struct {
		User                  string `yaml:"user"`
		Password              string `yaml:"password"`
		Host                  string `yaml:"host"`
		Database              string `yaml:"database"`
		MaxOpenConnections    int    `yaml:"max-open-connections"`
		MaxConnectionLifetime int    `yaml:"connection-lifetime"`
	} `yaml:"datasource"`
}

func LoadConfiguration(configPath string) {
	log.Printf("Loading configuration from %s", configPath)
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		log.Panicf("Error occured while trying to read configuration file: %s", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		log.Panicf("Error occured while trying to decode configuration values: %s", err)
	}

	ApplicationConfig = config
}
