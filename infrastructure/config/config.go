package config

import "flag"

type ServerConfig struct {
	Port string
	IP string
	Domain string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowedHeaders []string
	AllowedMethods []string
	ExposedHeaders []string
}

var GlobalServerConfig = ServerConfig{}

var GlobalCORSConfig = CORSConfig{
	AllowedOrigins: []string{"http://localhost", "https://softree.group", "http://localhost:3000"},
	AllowedHeaders: []string{"If-Modified-Since", "Cache-Control", "Content-Type", "Range"},
	AllowedMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
	ExposedHeaders: []string{"Content-Length", "Content-Range"},
}

func InitFlags() {
	flag.StringVar(&GlobalServerConfig.Port, "port", "8000", "-port 8000")
	flag.StringVar(&GlobalServerConfig.IP, "ip", "0.0.0.0", "-ip 127.0.0.1")
	flag.StringVar(&GlobalServerConfig.Domain, "domain", "localhost", "-domain http://localhost")

	flag.Parse()
}