package main

import (
	"net/http"
	"server/Handlers"
	"server/Handlers/Authorization"
	"server/Handlers/UserInteraction"
)

func main() {
	http.HandleFunc("/", Handlers.MainOrSignup)
	http.HandleFunc("/api/signin", Authorization.Login)
	http.HandleFunc("/api/signup", Authorization.Signup)
	http.HandleFunc("/api/logout", Authorization.Logout)
	http.HandleFunc("/api/user-data", UserInteraction.UserData)
	http.HandleFunc("/api/rates", UserInteraction.Rates)
	http.HandleFunc("/api/user", UserInteraction.UpdateUser)

	http.ListenAndServe(":8000", nil)
}
