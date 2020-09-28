package Handlers

import (
	"fmt"
	"net/http"
)

func MainOrSignup(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		fmt.Fprintf(w, "<h1>hello from root</h1>")
	} else {
		http.Redirect(w, r, signup, http.StatusUnauthorized)
	}
}
