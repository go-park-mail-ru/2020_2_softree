package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/config"
	"server/src/infrastructure/corsInteraction"
	"server/src/infrastructure/log"
	"server/src/infrastructure/persistence"
	"server/src/interfaces/authorization"
	"server/src/interfaces/profile"
	"server/src/interfaces/rates"
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

func createAuthenticate() *authorization.Authentication {
	memoryRepo := persistence.NewUserRepository()
	memoryAuth := auth.NewMemAuth()

	servicesDB := application.NewUserApp(memoryRepo)
	servicesAuth := application.NewUserAuth(memoryAuth)
	servicesLog := log.NewLogrusLogger()

	return authorization.NewAuthenticate(*servicesDB, *servicesAuth, servicesLog)
}

func createProfile() *profile.Profile {
	memoryRepo := persistence.NewUserRepository()
	memoryAuth := auth.NewMemAuth()

	servicesDB := application.NewUserApp(memoryRepo)
	servicesAuth := application.NewUserAuth(memoryAuth)
	servicesLog := log.NewLogrusLogger()

	return profile.NewProfile(*servicesDB, *servicesAuth, servicesLog)
}

func createRates() *rates.Rates {
	memoryRepo := persistence.NewRateRepository()
	memoryAuth := auth.NewMemAuth()

	servicesDB := application.NewRateApp(memoryRepo)
	servicesAuth := application.NewUserAuth(memoryAuth)
	servicesLog := log.NewLogrusLogger()

	return rates.NewRates(*servicesDB, *servicesAuth, servicesLog)
}

func main() {
	userAuthenticate := createAuthenticate()
	userProfile := createProfile()
	rateRates := createRates()

	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/signin", userAuthenticate.Login).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/signup", userAuthenticate.Signup).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/auth", userAuthenticate.Auth).
		Methods("GET", "OPTIONS")

	r.HandleFunc("/logout", userAuthenticate.Logout).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/rates", rateRates.GetRates).
		Methods("GET", "OPTIONS")

	r.HandleFunc("/user", userProfile.Auth(userProfile.UpdateUser)).
		Methods("PATCH", "OPTIONS")

	r.HandleFunc("/change-password", userProfile.Auth(userProfile.UpdateUser)).
		Methods("PATCH", "OPTIONS")

	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	server.ListenAndServe()
}
