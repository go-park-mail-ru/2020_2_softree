package authorization

import (
	"encoding/json"
	"errors"
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

	user, err = a.userApp.GetUserByLogin(user.Email, user.Password)
	errs := a.checkGetUserByLoginErrors(err)

	if errs.NotEmpty {
		a.log.Print(err)
		res, err := json.Marshal(errs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
		return
	}

	var cookie http.Cookie
	if cookie, err = a.auth.CreateAuth(user.ID); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (a *Authentication) checkGetUserByLoginErrors(err error) (errs entity.ErrorJSON) {
	if err == errors.New("user does not exist") {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "такой пользователь не существует")
	}
	if err == errors.New("wrong password") {
		errs.NotEmpty = true
		errs.Password = append(errs.Password, "неправильный пароль")
	}
	if err != nil {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, err.Error())
	}

	return
}
