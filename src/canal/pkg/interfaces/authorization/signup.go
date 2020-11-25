package authorization

import (
	"encoding/json"
	"net/http"
	session "server/src/authorization/pkg/session/gen"
	"server/src/canal/pkg/domain/entity"
	"server/src/canal/pkg/infrastructure/security"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Signup(w http.ResponseWriter, r *http.Request) {
	var user *profile.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Signup",
			"action":   "Decode",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
			"function": "Signup",
			"action":   "CheckExistence",
			"user":     user,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if check.Existence {
		errs.NonFieldError = append(errs.NonFieldError, "пользователь с таким email'ом уже существует")
		a.createServerError(&errs, w)
		return
	}

	if user.Password, err = security.MakeShieldedPassword(user.Password); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Signup",
			"action":   "MakeShieldedPassword",
			"user":     user,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	public, err := a.profile.SaveUser(r.Context(), user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Signup",
			"action":   "SaveUser",
			"user":     user,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cookie http.Cookie
	if cookie, err = CreateCookie(); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Signup",
			"action":   "CreateCookie",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = a.auth.Create(r.Context(), &session.Session{Id: public.Id, SessionId: cookie.Value}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Signup",
			"action":   "Create auth",
			"session":  session.Session{Id: user.Id, SessionId: cookie.Value},
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}
