package pureArchAuth

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

func (a *Authenticate) Signup(w http.ResponseWriter, r *http.Request) {
	var user *entity.User
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var errors *jsonRealisation.ErrorJSON
	user, errors = a.userApp.SaveUser(user)

	if errors.NotEmpty {
		res, err := json.Marshal(errors)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
	}

	cookie, err := a.cookie.CreateCookie()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &cookie)
	if err := a.auth.CreateAuth(user.ID, &cookie); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
