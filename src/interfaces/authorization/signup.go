package authorization

import (
	"encoding/json"
	"net/http"
	session "server/src/authorizationService/session/gen"
	"server/src/domain/entity"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Signup(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a.sanitizer.Sanitize(user.Email)
	a.sanitizer.Sanitize(user.Password)

	errs := user.Validate()
	if errs.NotEmpty {
		a.createServerError(&errs, w)
		return
	}

	var exist bool
	if exist, err = a.userApp.CheckExistence(user.Email); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if exist {
		errs.NonFieldError = append(errs.NonFieldError, "пользователь с таким email'ом уже существует")
		a.createServerError(&errs, w)
		return
	}

	user, err = a.userApp.SaveUser(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cookie http.Cookie
	if cookie, err = CreateCookie(); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = a.auth.Create(r.Context(), &session.Session{Id: user.ID, SessionId: cookie.Value}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}
