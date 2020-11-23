package authorization

import (
	"context"
	"encoding/json"
	"net/http"
	session "server/src/authorization/session/gen"
	"server/src/canal/domain/entity"

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

		ctx := context.WithValue(r.Context(), "id", id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authentication) Authenticate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(*session.UserID)

	var user entity.User
	var err error
	if user, err = a.userApp.GetUserById(id.Id); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusBadRequest}).Error(err)
	}
}
