package authorization

import (
	"context"
	"encoding/json"
	"net/http"
	session "server/authorization/pkg/session/gen"
	profile "server/profile/pkg/profile/gen"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := a.auth.Check(r.Context(), &session.SessionID{SessionId: cookie.Value})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"status": http.StatusBadRequest}).Error(err)

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// линтер ругается если используем базовые типы в Value контекста
		// типа так безопаснее разграничивать
		type key string

		const idKey key = "id"

		ctx := context.WithValue(r.Context(), idKey, id.Id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authentication) Authenticate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	var user *profile.PublicUser
	var err error
	if user, err = a.profile.GetUserById(r.Context(), &profile.UserID{Id: id}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusBadRequest,
			"function": "Authenticate",
			"UserID":   id,
			"action":   "GetUserById",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "Authenticate",
			"UserID":   id,
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "Authenticate",
			"UserID":   id,
			"action":   "Write",
		}).Error(err)
	}
}
