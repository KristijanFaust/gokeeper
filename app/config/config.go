package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	*Profile        `yaml:"profile"`
	*Server         `yaml:"server"`
	*Datasource     `yaml:"datasource"`
	*Authentication `yaml:"authentication"`
}

type Profile struct {
	Production bool `yaml:"production"`
}

type Server struct {
	Hostname string `yaml:"hostname"`
	Port     string `yaml:"port"`
}

type Datasource struct {
	User                  string `yaml:"user"`
	Password              string `yaml:"password"`
	Host                  string `yaml:"host"`
	Database              string `yaml:"database"`
	MaxOpenConnections    int    `yaml:"max-open-connections"`
	MaxConnectionLifetime int    `yaml:"connection-lifetime"`
}

type Authentication struct {
	Issuer               string `yaml:"issuer"`
	JwtSigningKey        string `yaml:"jwt-signing-key"`
	JwtDurationInMinutes int    `yaml:"jwt-duration-in-minutes"`
}

func LoadConfiguration(configPath string) *Config {
	log.Printf("Loading configuration from %s", configPath)
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		log.Panicf("Error occured while trying to read configuration file: %s", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(config); err != nil {
		log.Panicf("Error occured while trying to decode configuration values: %s", err)
	}

	return config
}
