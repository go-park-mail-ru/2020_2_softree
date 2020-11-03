package profile

import (
	"context"
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (p *Profile) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := p.auth.CheckAuth(cookie.Value)
		if err != nil {
			p.log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "id", id)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (p *Profile) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user, err = p.userApp.UpdateUser(id, user); err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(user.MakePublicUser())
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
}
