package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/domain/entity/rates"
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

	r.HandleFunc("/signin", login.Login).Methods("POST")
	r.HandleFunc("/signup", signup.Signup).Methods("POST")
	r.HandleFunc("/logout", logout.Logout)
	r.HandleFunc("/user-data", userInteraction.UserData).Methods("POST")
	r.HandleFunc("/rates", ratesInteraction.Rates).Methods("GET")
	r.HandleFunc("/rates/{id:([1-9]0?)+}", ratesInteraction.Rates).Methods("GET")
	r.HandleFunc("/user", userInteraction.UpdateUser).Methods("PUT", "PATCH")
	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go rates.StartTicker()

	log.Fatal(server.ListenAndServe())
}
