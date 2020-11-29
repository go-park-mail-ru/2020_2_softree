package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"server/canal/pkg/infrastructure/CORS"
	"server/canal/pkg/interfaces/authorization"
	"server/canal/pkg/interfaces/profile"
	"server/canal/pkg/interfaces/rates"
)

func NewRouter(userAuthenticate *authorization.Authentication, userProfile *profile.Profile, rateRates *rates.Rates) http.Handler {
	router := mux.NewRouter()
	r := router.PathPrefix("/api").Subrouter()

	r.HandleFunc("/sessions", userAuthenticate.Login).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/sessions", userAuthenticate.Logout).
		Methods(http.MethodDelete, http.MethodOptions)

	r.HandleFunc("/sessions", userAuthenticate.Auth(userAuthenticate.Authenticate)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/rates", rateRates.GetRates).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/rates/{title}", rateRates.GetURLRate).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users", userAuthenticate.Signup).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/users", userAuthenticate.Auth(userProfile.UpdateUserAvatar)).
		Methods(http.MethodPut, http.MethodOptions)

	r.HandleFunc("/users", userAuthenticate.Auth(userProfile.GetUser)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/users/change-password", userAuthenticate.Auth(userProfile.UpdateUserPassword)).
		Methods(http.MethodPut, http.MethodOptions)

	r.HandleFunc("/watchers", userAuthenticate.Auth(userProfile.GetUserWatchlist)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/markets", rateRates.GetMarkets).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/accounts", userAuthenticate.Auth(userProfile.GetWallets)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/accounts", userAuthenticate.Auth(userProfile.SetWallet)).
		Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/transactions", userAuthenticate.Auth(userProfile.GetTransactions)).
		Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/transactions", userAuthenticate.Auth(userProfile.SetTransaction)).
		Methods(http.MethodPost, http.MethodOptions)

	r.Use(CORS.CORSMiddleware())

	return r
}
