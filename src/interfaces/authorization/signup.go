package authorization

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (a *Authentication) Signup(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.sanitizer.Sanitize(user.Email)
	a.sanitizer.Sanitize(user.Password)

	errs := user.Validate()
	if errs.NotEmpty {
		a.createInternalServerError(&errs, w)
		return
	}

	var exist bool
	if exist, err = a.userApp.CheckExistence(user.Email); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exist {
		errs.NonFieldError = append(errs.NonFieldError, "пользователь с таким email'ом уже существует")
		a.createInternalServerError(&errs, w)
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

func (a *Authentication) createInternalServerError(errors *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		a.log.Print(err)
	}
}
