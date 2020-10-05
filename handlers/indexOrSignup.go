package handlers

import (
	"net/http"
	"server/handlers/authorization/utils"
)

func MainOrSignup(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		// find user in db session table to write his data in json
		http.Redirect(w, r, utils.RootPage, http.StatusOK)
	} else {
		http.Redirect(w, r, utils.SignupPage, http.StatusUnauthorized)
	}
}
