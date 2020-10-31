package pureArchAuth

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

func (a *Authenticate) Signup(w http.ResponseWriter, r *http.Request) {
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

	user, err = a.userApp.SaveUser(user)
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cookie http.Cookie
	if cookie, err = a.auth.CreateAuth(user.ID); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}

func (a *Authenticate) createInternalServerError(errors *jsonRealisation.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
