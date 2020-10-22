package pureArchAuth

import (
	"net/http"
	"server/src/infrastructure/auth"
	"time"
)

func (a *Authenticate) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	if err = a.auth.DeleteAuth(&auth.AccessDetails{Value: cookie.Value}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie, err := a.cookie.CreateCookie()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusFound)
}