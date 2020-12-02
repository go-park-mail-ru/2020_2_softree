package authorization

import (
	"encoding/json"
	"net/http"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"

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

		a.recordHitMetric(http.StatusInternalServerError)
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

		a.recordHitMetric(http.StatusInternalServerError)
		return
	}
	if !check.Existence {
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
		a.createServerError(&errs, w)
		return
	}

	public, err := a.profile.GetUserByLogin(r.Context(), user)
	if err != nil {
		if errs = a.checkGetUserByLoginErrors(err); errs.NotEmpty {
			a.createServerError(&errs, w)
			return
		}
	}

	cookie := CreateCookie()
	var sess *session.Session
	var userId = &session.UserID{Id: public.Id}
	if sess, err = a.auth.Create(r.Context(), userId); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "Create auth",
			"session":  &session.UserID{Id: public.Id},
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		a.recordHitMetric(http.StatusInternalServerError)
		return
	}
	cookie.Value = sess.SessionId
	http.SetCookie(w, &cookie)

	res, err := json.Marshal(public)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Login",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		a.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	a.recordHitMetric(http.StatusOK)

	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "Login",
			"action":   "Write",
		}).Error(err)
	}
}
