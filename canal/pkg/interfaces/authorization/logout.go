package authorization

import (
	"net/http"
	session "server/authorization/pkg/session/gen"
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusUnauthorized)
		a.recordHitMetric(http.StatusUnauthorized)
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

		a.recordHitMetric(http.StatusInternalServerError)
		return
	}

	newCookie := CreateCookie()
	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""
	http.SetCookie(w, &newCookie)
	w.WriteHeader(http.StatusOK)

	a.recordHitMetric(http.StatusOK)
}
