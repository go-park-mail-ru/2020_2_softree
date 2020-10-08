package config

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port   string `yaml:"port"`
	IP     string `yaml:"ip"`
	Domain string `yaml:"domain"`
	Secure bool   `yaml:"secure"`
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

func ParseConfig() {
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &GlobalServerConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func InitFlags() {
	flag.StringVar(&GlobalServerConfig.Port, "port", "8000", "-port 8000")
	flag.StringVar(&GlobalServerConfig.IP, "ip", "0.0.0.0", "-ip 127.0.0.1")
	flag.StringVar(&GlobalServerConfig.Domain, "domain", "localhost", "-domain http://localhost")
	flag.BoolVar(&GlobalServerConfig.Secure, "secure", false, "-secure true")

	flag.StringVar(&configPath, "f", "", "-f path to config file")

	flag.Parse()

	if configPath != "" {
		ParseConfig()
	}
}
