package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/src/domain/entity/rates"
	"server/src/handlers/authorization/auth"
	"server/src/handlers/authorization/login"
	"server/src/handlers/authorization/logout"
	"server/src/handlers/authorization/signup"
	"server/src/handlers/ratesInteraction"
	"server/src/handlers/userInteraction"
	"server/src/infrastructure/config"
	"server/src/infrastructure/corsInteraction"
	"time"
)

func main() {
	config.InitFlags()
	go rates.StartTicker()

	router := mux.NewRouter()
	r := router.PathPrefix("").Subrouter()

	r.HandleFunc("/signin", login.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/signup", signup.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/auth", auth.Authentication).Methods("GET", "OPTIONS")
	r.HandleFunc("/logout", logout.Logout).Methods("POST", "OPTIONS")
	r.HandleFunc("/rates", ratesInteraction.Rates).Methods("GET", "OPTIONS")
	r.HandleFunc("/user", userInteraction.UpdateUser).Methods("PUT", "PATCH", "OPTIONS")
	r.HandleFunc("/change-password", userInteraction.UpdatePassword).Methods("PATCH", "OPTIONS")
	r.Use(corsInteraction.CORSMiddleware())

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port),
		Handler:      r,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
