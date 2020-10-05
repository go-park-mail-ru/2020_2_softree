package config

import "flag"

type Flags struct {
	Port string
	IP string
	CORSDomain string
}

var Options = new(Flags)

func InitFlags() {
	flag.StringVar(&Options.Port, "port", "8000", "-port 8000")
	flag.StringVar(&Options.IP, "ip", "0.0.0.0", "-ip 127.0.0.1")
	flag.StringVar(&Options.CORSDomain, "cors_domain", "http://localhost", "-cors_domain http://example.com")

	flag.Parse()
}