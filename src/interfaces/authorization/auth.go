package authorization

import (
	"context"
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (a *Authentication) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := a.auth.CheckAuth(cookie.Value)
		if err != nil {
			a.log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "id", id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (a *Authentication) Authenticate(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var user entity.User
	var err error
	if user, err = a.userApp.GetUserById(id); err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		a.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		a.log.Print(err)
	}
}
