package main

import (
	"log"
	"net/http"
	"server/handlers"
	"server/handlers/authorization/login"
	"server/handlers/authorization/logout"
	"server/handlers/authorization/signup"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
)

func main() {
	http.HandleFunc("/", handlers.MainOrSignup)
	http.HandleFunc("/api/signin", login.Login)
	http.HandleFunc("/api/signup", signup.Signup)
	http.HandleFunc("/api/logout", logout.Logout)
	http.HandleFunc("/api/user-data", userInteraction.UserData)
	http.HandleFunc("/api/rates", ratesInteraction.Rates)
	http.HandleFunc("/api/user", userInteraction.UpdateUser)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
