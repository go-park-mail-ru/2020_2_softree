package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/domain/entity/rates"
	"server/handlers/authorization/auth"
	"server/handlers/authorization/login"
	"server/handlers/authorization/logout"
	"server/handlers/authorization/signup"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
	"server/infrastructure/config"
	"server/infrastructure/corsInteraction"
	"time"
)

func main() {
	config.InitFlags()

	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/signin", login.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", signup.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth", auth.Authentication).Methods("GET", "OPTIONS")
	r.HandleFunc("/logout", logout.Logout)
	r.HandleFunc("/user-data", userInteraction.UserData).Methods("POST", "OPTIONS")
	r.HandleFunc("/rates", ratesInteraction.Rates).Methods("GET", "OPTIONS")
	// r.HandleFunc("/rates/{id:([1-9]0?)+}", ratesInteraction.Rates).Methods("GET")
	r.HandleFunc("/user", userInteraction.UpdateUser).Methods("PUT", "PATCH", "OPTIONS")
	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go rates.StartTicker()

	log.Fatal(server.ListenAndServe())
}
