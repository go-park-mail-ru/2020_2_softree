package authorization

import (
	"encoding/json"
	"net/http"
	session "server/src/authorization/session/gen"
	"server/src/canal/domain/entity"
	profile "server/src/profile/profile/gen"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Login(w http.ResponseWriter, r *http.Request) {
	var user *profile.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "Decode",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	a.sanitizer.Sanitize(user.Email)
	a.sanitizer.Sanitize(user.Password)

	var errs entity.ErrorJSON
	if errs = a.validate(user); errs.NotEmpty {
		a.createServerError(&errs, w)
		return
	}

	var check *profile.Check
	if check, err = a.profile.CheckExistence(r.Context(), user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "CheckExistence",
			"user":     user,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check.Existence {
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
		a.createServerError(&errs, w)
		return
	}

	public, err := a.profile.GetUserByLogin(r.Context(), user)
	if errs = a.checkGetUserByLoginErrors(err); errs.NotEmpty {
		a.createServerError(&errs, w)
		return
	}

	var cookie http.Cookie
	if cookie, err = CreateCookie(); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "CreateCookie",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = a.auth.Create(r.Context(), &session.Session{Id: public.ID, SessionId: cookie.Value}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "Create auth",
			"session":  session.Session{Id: public.ID, SessionId: cookie.Value},
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	res, err := json.Marshal(public)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "Login",
			"action":   "Write",
		}).Error(err)
	}
}
