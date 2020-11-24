package authorization

import (
	"net/http"
	session "server/src/authorization/pkg/session/gen"
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if _, err = a.auth.Delete(r.Context(), &session.SessionID{SessionId: cookie.Value}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":     http.StatusInternalServerError,
			"function":   "Logout",
			"action":     "Delete auth",
			"session_id": session.SessionID{SessionId: cookie.Value},
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie, err := CreateCookie()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Logout",
			"action":   "CreateCookie",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusOK)
}
