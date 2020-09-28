package main

import (
	"net/http"
	"server/Handlers"
)

func main() {
	http.HandleFunc("/", Handlers.MainOrSignup)
	http.HandleFunc("/api/signin", Handlers.Login)
	http.HandleFunc("/api/signup", Handlers.Signup)
	http.HandleFunc("/api/logout", Handlers.Logout)
	http.HandleFunc("/api/user-data", Handlers.UserData)
	http.HandleFunc("/api/rates", Handlers.Rates)
	http.HandleFunc("/api/user", Handlers.UpdateUser)

	http.ListenAndServe(":8000", nil)
}
