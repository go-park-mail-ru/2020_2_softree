package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/infrastructure/config"
	"server/src/interfaces/router"
	"time"
)

func initFlags() {
	var helpFlag bool

	flag.StringVar(&config.GlobalServerConfig.Port, "p", "8000", "-p set port to listen")
	flag.StringVar(&config.GlobalServerConfig.IP, "ip", "127.0.0.1", "-ip set ip addr")
	flag.StringVar(&config.GlobalServerConfig.Domain, "d", "localhost", "-d set domain name")
	flag.BoolVar(&config.GlobalServerConfig.Secure, "s", false, "-s set cookie HTTPS only")
	flag.StringVar(&config.GlobalServerConfig.ConfigFile, "f", "", "-f path to config file")
	flag.StringVar(&config.GlobalServerConfig.LogLevel, "ll", "Info", "-ll set log level")
	flag.StringVar(&config.GlobalServerConfig.LogFile, "lf", "", "-lf set log file")
	flag.BoolVar(&helpFlag, "h", false, "-h get usage message")

	// Databases configs
	flag.StringVar(&config.SessionDatabaseConfig.AddressSessions, "rp",
		"redis://user:@localhost:6379/1", "set redis session database addr")
	flag.StringVar(&config.SessionDatabaseConfig.AddressDayCurrency, "rc",
		"redis://user:@localhost:6379/2", "set redis day currency database addr")

	flag.StringVar(&config.UserDatabaseConfig.User, "pu", "app_user", "User DB user")
	flag.StringVar(&config.UserDatabaseConfig.Password, "pp", "NeverGonnaGiveYouUp", "User DB password")
	flag.StringVar(&config.UserDatabaseConfig.Host, "ph", "localhost", "User DB port")
	flag.StringVar(&config.UserDatabaseConfig.Schema, "pd", "users", "User DB schema")

	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if config.GlobalServerConfig.ConfigFile != "" {
		if err := config.ParseConfig(); err != nil {
			fmt.Fprint(os.Stderr, "Error during parsing config", err)
			os.Exit(1)
		}
		return
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	initFlags()
}

func main() {
	userAuthenticate, userProfile, rateRates, err := router.CreateAppStructs()
	if err != nil {
		fmt.Print(err)
		return
	}

	go rateRates.GetRatesFromApi()

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      router.NewRouter(userAuthenticate, userProfile, rateRates),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	server.ListenAndServe()
}
