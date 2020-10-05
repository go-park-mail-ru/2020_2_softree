package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	//"server/handlers"
	"server/handlers/authorization/login"
	"server/handlers/authorization/logout"
	"server/handlers/authorization/signup"
	"server/handlers/ratesInteraction"
	"server/handlers/userInteraction"
	"server/infrastructure/config"
)

func enableCORS(cfg *config.CORSConfig, handler http.Handler) http.Handler {
	var (
		allowedOrigins = handlers.AllowedOrigins(cfg.AllowedOrigins)
		allowedHeaders = handlers.AllowedHeaders(cfg.AllowedHeaders)
		exposedHeaders = handlers.ExposedHeaders(cfg.ExposedHeaders)
		allowedMethods = handlers.AllowedMethods(cfg.AllowedMethods)
		credentials = handlers.AllowCredentials()
	)

	return handlers.CORS(allowedOrigins, allowedHeaders, exposedHeaders, allowedMethods, credentials)(handler)
}

func CORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return enableCORS(&config.GlobalCORSConfig, next)
	}
}

func main() {
	config.InitFlags()

	r := mux.NewRouter()

	r.HandleFunc("/api/signin", login.Login)
	r.HandleFunc("/api/signup", signup.Signup)
	r.HandleFunc("/api/logout", logout.Logout)
	r.HandleFunc("/api/user-data", userInteraction.UserData)
	r.HandleFunc("/api/rates", ratesInteraction.Rates)
	r.HandleFunc("/api/user", userInteraction.UpdateUser)
	r.Use(CORSMiddleware())

	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.GlobalServerConfig.IP, config.GlobalServerConfig.Port), r)

	if err != nil {
		log.Fatal(err)
	}
}
