package main

import (
	"net/http"
	"server/handlers"
	"server/handlers/authorization"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
)

func main() {
	http.HandleFunc("/", handlers.MainOrSignup)
	http.HandleFunc("/api/signin", authorization.Login)
	http.HandleFunc("/api/signup", authorization.Signup)
	http.HandleFunc("/api/logout", authorization.Logout)
	http.HandleFunc("/api/user-data", userInteraction.UserData)
	http.HandleFunc("/api/rates", ratesInteraction.Rates)
	http.HandleFunc("/api/user", userInteraction.UpdateUser)

	http.ListenAndServe(":8000", nil)
}
