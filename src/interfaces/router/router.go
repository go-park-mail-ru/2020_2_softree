package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"server/src/infrastructure/corsInteraction"
	"server/src/interfaces/authorization"
	"server/src/interfaces/profile"
	"server/src/interfaces/rates"
)

func NewRouter(userAuthenticate *authorization.Authentication, userProfile *profile.Profile, rateRates *rates.Rates) http.Handler {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/sessions", userAuthenticate.Login).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/sessions", userAuthenticate.Logout).
		Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/sessions", userAuthenticate.Auth).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/rates", rateRates.GetRates).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/rates/{title}/", rateRates.GetURLRate).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users", userAuthenticate.Signup).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/users", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodPut, http.MethodOptions)

	r.HandleFunc("/users", userProfile.Auth(userProfile.GetUser)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/watchers", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users/change-password", userProfile.Auth(userProfile.UpdateUser)).
		Methods(http.MethodPut, http.MethodOptions)

	r.HandleFunc("/markets", rateRates.GetMarkets).
		Methods(http.MethodGet, http.MethodOptions)

	r.Use(corsInteraction.CORSMiddleware())

	return r
}
