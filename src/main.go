package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/domain/entity/rates"
	"server/src/infrastructure/config"
	"server/src/infrastructure/corsInteraction"
	"server/src/infrastructure/log"
	"server/src/interfaces/authorization/auth"
	"server/src/interfaces/authorization/login"
	"server/src/interfaces/authorization/logout"
	"server/src/interfaces/authorization/signup"
	"server/src/interfaces/ratesInteraction"
	"server/src/interfaces/profile"
	"time"

	"github.com/gorilla/mux"
)

func initFlags() {
	var helpFlag bool

	flag.StringVar(&config.GlobalServerConfig.Port, "p", "8000", "-p set port to listen")
	flag.StringVar(&config.GlobalServerConfig.IP, "ip", "127.0.0.1", "-ip set ip addr")
	flag.StringVar(&config.GlobalServerConfig.Domain, "d", "", "-d set domain name")
	flag.BoolVar(&config.GlobalServerConfig.Secure, "s", false, "-s set cookie HTTPS only")
	flag.StringVar(&config.GlobalServerConfig.ConfigFile, "f", "", "-f path to config file")
	flag.StringVar(&config.GlobalServerConfig.LogLevel, "ll", "Info", "-ll set log level")
	flag.StringVar(&config.GlobalServerConfig.LogFile, "lf", "", "-lf set log file")
	flag.BoolVar(&helpFlag, "h", false, "-h get usage message")

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
	go rates.StartTicker()
	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/signin", login.Login).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/signup", signup.Signup).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/auth", auth.Authentication).
		Methods("GET", "OPTIONS")

	r.HandleFunc("/logout", logout.Logout).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/rates", ratesInteraction.Rates).
		Methods("GET", "OPTIONS")

	r.HandleFunc("/user", profile.UpdateUserPartly).
		Methods("PATCH", "OPTIONS")

	r.HandleFunc("/change-password", profile.UpdatePassword).
		Methods("PATCH", "OPTIONS")

	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.GlobalLogger.Error(server.ListenAndServe())
}
