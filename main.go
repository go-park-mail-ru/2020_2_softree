package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"server/handlers/authorization/login"
	"server/handlers/authorization/logout"
	"server/handlers/authorization/signup"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
	"server/infrastructure/config"
	"server/infrastructure/corsInteraction"
)

func main() {
	config.InitFlags()

	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/signin", login.Login)
	r.HandleFunc("/signup", signup.Signup)
	r.HandleFunc("/logout", logout.Logout)
	r.HandleFunc("/user-data", userInteraction.UserData)
	r.HandleFunc("/rates", ratesInteraction.Rates)
	r.HandleFunc("/user", userInteraction.UpdateUser)
	r.Use(corsInteraction.CORSMiddleware())

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
