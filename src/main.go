package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"server/src/domain/entity/rates"
	"server/src/handlers/authorization/auth"
	"server/src/handlers/authorization/login"
	"server/src/handlers/authorization/logout"
	"server/src/handlers/authorization/signup"
	"server/src/handlers/ratesInteraction"
	"server/src/handlers/userInteraction"
	"server/src/infrastructure/config"
	"server/src/infrastructure/corsInteraction"
	"server/src/infrastructure/log"
	"time"

	"github.com/gorilla/mux"
)

func initFlags() {
	var helpFlag bool

	flag.StringVar(&config.GlobalServerConfig.Port, "p", "", "-p set port to listen")
	flag.StringVar(&config.GlobalServerConfig.IP, "ip", "", "-ip set ip addr")
	flag.StringVar(&config.GlobalServerConfig.Domain, "d", "", "-d set domain")
	flag.BoolVar(&config.GlobalServerConfig.Secure, "s", true, "-s set CORS")
	flag.StringVar(&config.GlobalServerConfig.ConfigFile, "f", "", "-f path to config file")
	flag.StringVar(&config.GlobalServerConfig.LogLevel, "ll", "info", "-ll set log level")
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

	if config.GlobalServerConfig.IP == "" ||
		config.GlobalServerConfig.Port == "" ||
		config.GlobalServerConfig.Domain == "" {
		fmt.Fprint(os.Stderr, "Need to explicit set server ip:port and domain")
		flag.Usage()
		os.Exit(1)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	initFlags()

	if err := log.ConfigureLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot inizialize logger %v\n", err)
		os.Exit(1)
	}
}

func main() {
	go rates.StartTicker()
	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/signin", login.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", signup.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth", auth.Authentication).Methods("GET", "OPTIONS")
	r.HandleFunc("/logout", logout.Logout).Methods("POST", "OPTIONS")
	r.HandleFunc("/rates", ratesInteraction.Rates).Methods("GET", "OPTIONS")
	r.HandleFunc("/user", userInteraction.UpdateUserPartly).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/change-password", userInteraction.UpdatePassword).Methods("PATCH", "OPTIONS")
	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.GlobalLogger.Error(server.ListenAndServe())
}
