package userInteraction

import (
	"encoding/json"
	"net/http"
	"server/Domain/Entity"
	"server/Handlers/Authorization"
)

const (
	testID    = 123
	testEmail = "hound@psina.ru"
)

func UserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	var u Entity.User
	u.ID = testID
	u.Email = testEmail
	result, e := json.Marshal(u)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if logged {
		w.Write(result)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, Authorization.SignupPage, http.StatusUnauthorized)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
