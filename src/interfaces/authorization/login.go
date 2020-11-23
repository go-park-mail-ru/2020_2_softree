package authorization

import (
	"encoding/json"
	"net/http"
	session "server/src/authorizationService/session/gen"
	"server/src/domain/entity"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

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

	if !exist {
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
		a.createServerError(&errs, w)
		return
	}

	user, err = a.userApp.GetUserByLogin(user.Email, user.Password)
	if errs = a.checkGetUserByLoginErrors(err); errs.NotEmpty {
		a.createServerError(&errs, w)
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

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
	}
}
