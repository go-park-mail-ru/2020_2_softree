package main

import (
	"fmt"
	"log"
	"net/http"
	"server/handlers"
	"server/handlers/authorization/login"
	"server/handlers/authorization/logout"
	"server/handlers/authorization/signup"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
	"server/infrastructure/config"
)

func main() {
	config.InitFlags()

	http.HandleFunc("/", handlers.MainOrSignup)
	http.HandleFunc("/api/signin", login.Login)
	http.HandleFunc("/api/signup", signup.Signup)
	http.HandleFunc("/api/logout", logout.Logout)
	http.HandleFunc("/api/user-data", userInteraction.UserData)
	http.HandleFunc("/api/rates", ratesInteraction.Rates)
	http.HandleFunc("/api/user", userInteraction.UpdateUser)

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Options.IP, config.Options.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
