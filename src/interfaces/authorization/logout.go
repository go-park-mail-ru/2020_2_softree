package authorization

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = a.auth.DeleteAuth(cookie.Value); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie, err := a.auth.CreateCookie()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusOK)
}
