package userInteraction

import (
	"encoding/json"
	"net/http"
	"server/domain/entity"
)

const (
	testID    = 123
	testEmail = "hound@psina.ru"
)

func UserData(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var u entity.User
		u.ID = testID
		u.Email = testEmail
		result, e := json.Marshal(u)
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
