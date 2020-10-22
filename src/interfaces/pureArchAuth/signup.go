package pureArchAuth

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

func (a *Authenticate) Signup(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	errors := user.Validate("signup")
	if errors.NotEmpty {
		createInternalServerError(&errors, w)
		return
	}
	user, _ = a.userApp.SaveUser(user)

	cookie, err := a.cookie.CreateCookie()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	if err := a.auth.CreateAuth(user.ID, cookie.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func createInternalServerError(errors *jsonRealisation.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}