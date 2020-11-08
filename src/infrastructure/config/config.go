package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port       string `yaml:"port"`
	IP         string `yaml:"ip"`
	Domain     string `yaml:"domain"`
	Secure     bool   `yaml:"secure"`
	LogLevel   string `yaml:"logLevel"`
	LogFile    string `yaml:"logFile"`
	ConfigFile string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
	ExposedHeaders []string
}

var GlobalServerConfig = ServerConfig{}

var GlobalCORSConfig = CORSConfig{
	AllowedOrigins: []string{"http://localhost", "https://softree.group", "http://localhost:3000", "http://self.ru"},
	AllowedHeaders: []string{"If-Modified-Since", "Cache-Control", "Content-Type", "Range"},
	AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
	ExposedHeaders: []string{"Content-Length", "Content-Range"},
}

type RedisConfig struct {
	AddressSessions    string
	AddressDayCurrency string
}

var SessionDatabaseConfig = RedisConfig{}

type RateBDConfig struct {
	User     string
	Password string
	Port     string
	Schema   string
}

var RateDatabaseConfig = RateBDConfig{}

type UserBDConfig struct {
	User     string
	Password string
	Port     string
	Schema   string
}

var UserDatabaseConfig = UserBDConfig{}

func ParseConfig() error {
	yamlFile, err := ioutil.ReadFile(GlobalServerConfig.ConfigFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &GlobalServerConfig)
	if err != nil {
		return err
	}
	return nil
}
