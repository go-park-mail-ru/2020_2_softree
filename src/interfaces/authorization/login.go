package authorization

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (a *Authentication) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errors := user.Validate()
	if errors.NotEmpty {
		a.createInternalServerError(&errors, w)
		return
	}

	if user, err = a.userApp.GetUserByLogin(user.Email, user.Password); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cookie http.Cookie
	if cookie, err = a.auth.CreateAuth(user.ID); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
