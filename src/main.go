package main

import (
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"net/http"
	"os"
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/config"
	"server/src/infrastructure/corsInteraction"
	"server/src/infrastructure/financial"
	"server/src/infrastructure/log"
	"server/src/infrastructure/persistence"
	"server/src/infrastructure/userInteraction"
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

	// Databases configs
	flag.StringVar(&config.SessionDatabaseConfig.AddressSessions, "redisSessions",
		"redis://user:@localhost:6379/1", "set redis session database addr")
	flag.StringVar(&config.SessionDatabaseConfig.AddressDayCurrency, "redisDayCurrency",
		"redis://user:@localhost:6379/2", "set redis day currency database addr")

	flag.StringVar(&config.RateDatabaseConfig.User, "rate_db_user", "app_rates", "rate DB user")
	flag.StringVar(&config.RateDatabaseConfig.Password, "rate_db_password", "NeverGonnaGiveYouUp", "rate DB password")
	flag.StringVar(&config.RateDatabaseConfig.Host, "rate_db_host", "localhost", "rate DB port")
	flag.StringVar(&config.RateDatabaseConfig.Schema, "rate_db_schema", "rates", "rate DB schema")

	flag.StringVar(&config.UserDatabaseConfig.User, "user_db_user", "app_user", "User DB user")
	flag.StringVar(&config.UserDatabaseConfig.Password, "user_db_password", "NeverGonnaGiveYouUp", "User DB password")
	flag.StringVar(&config.UserDatabaseConfig.Host, "user_db_host", "localhost", "User DB port")
	flag.StringVar(&config.UserDatabaseConfig.Schema, "user_db_schema", "users", "User DB schema")

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

func createAuthenticate() (*authorization.Authentication, error) {
	dbRepo, err := userInteraction.NewUserDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressSessions)
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	servicesDB := application.NewUserApp(dbRepo)
	servicesAuth := application.NewUserAuth(dbAuth)
	servicesLog := log.NewLogrusLogger()

	return authorization.NewAuthenticate(*servicesDB, *servicesAuth, servicesLog), nil
}

func createProfile() (*profile.Profile, error) {
	dbRepo, err := userInteraction.NewUserDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressSessions)
	if err != nil {
		return nil, err
	}
	dbAuth := auth.NewSessionManager(connect)

	servicesDB := application.NewUserApp(dbRepo)
	servicesAuth := application.NewUserAuth(dbAuth)
	servicesLog := log.NewLogrusLogger()

	return profile.NewProfile(*servicesDB, *servicesAuth, servicesLog), nil
}

func createRates() (*rates.Rates, error) {
	dbRepo, err := persistence.NewRateDBManager()
	if err != nil {
		return nil, err
	}

	connect, err := redis.DialURL(config.SessionDatabaseConfig.AddressDayCurrency)
	if err != nil {
		return nil, err
	}
	dbAuth := financial.NewCurrencyManager(connect)

	servicesDB := application.NewRateApp(dbRepo, dbAuth)
	servicesLog := log.NewLogrusLogger()

	return rates.NewRates(*servicesDB, servicesLog), nil
}

func main() {
	userAuthenticate, err := createAuthenticate()
	if err != nil {
		fmt.Println("userAuthenticate", err)
		return
	}

	userProfile, err := createProfile()
	if err != nil {
		fmt.Println("userProfile", err)
		return
	}

	rateRates, err := createRates()
	if err != nil {
		fmt.Println("userProfile", err)
		return
	}

	go rateRates.GetRatesFromApi()

	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/sessions", userAuthenticate.Login).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/users", userAuthenticate.Signup).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/auth", userAuthenticate.Auth).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/logout", userAuthenticate.Logout).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/rates", rateRates.GetRates).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/rates/{title}", rateRates.GetURLRate).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodPut, http.MethodOptions)

	// #TODO
	r.HandleFunc("/users", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodGet, http.MethodOptions)

	// #TODO
	r.HandleFunc("/watchers", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users/change-password", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodPut, http.MethodOptions)

	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	server.ListenAndServe()
}
