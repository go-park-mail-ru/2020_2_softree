package config

import (
	"errors"
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port     string `yaml:"port"`
	IP       string `yaml:"ip"`
	Domain   string `yaml:"domain"`
	Secure   bool   `yaml:"secure"`
	LogLevel string `yaml:"logLevel"`
	LogFile  string `yaml:"logFile"`
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

var configPath string

func ParseConfig() error {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &GlobalServerConfig)
	if err != nil {
		return err
	}
	return nil
}

func InitFlags() error {
	flag.StringVar(&GlobalServerConfig.Port, "p", "", "-p set port to listen")
	flag.StringVar(&GlobalServerConfig.IP, "ip", "", "-ip set ip addr")
	flag.StringVar(&GlobalServerConfig.Domain, "d", "", "-d set domain")
	flag.BoolVar(&GlobalServerConfig.Secure, "s", true, "-s set CORS")
	flag.StringVar(&configPath, "f", "", "-f path to config file")
	flag.StringVar(&GlobalServerConfig.LogLevel, "ll", "info", "-ll set log level")
	flag.StringVar(&GlobalServerConfig.LogFile, "lf", "", "-lf set log file")

	flag.Parse()

	if configPath != "" {
		if err := ParseConfig(); err != nil {
			return err
		}
		return nil
	}

	if GlobalServerConfig.IP == "" ||
		GlobalServerConfig.Port == "" ||
		GlobalServerConfig.Domain == "" {
		flag.Usage()
		return errors.New("Need to explicity set ip:port and domain")
	}

	return nil
}
