package pureArchAuth

import (
	"encoding/json"
	"net/http"
)

func (a *Authenticate) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := a.auth.CheckAuth(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := a.userApp.GetUser(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
