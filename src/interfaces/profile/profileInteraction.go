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

	p.sanitizer.SanitizeBytes([]byte(user.Avatar))
	p.sanitizer.Sanitize(user.OldPassword)
	p.sanitizer.Sanitize(user.NewPassword)

	user, err = p.userApp.UpdateUser(id, user)
	if err != nil {
		if err.Error() == "wrong old password" {
			w.WriteHeader(http.StatusBadRequest)
			p.createOldPassError(w)
			return
		}
		p.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *Profile) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var user entity.User
	var err error
	if user, err = p.userApp.GetUserById(id); err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
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
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
}

func (p *Profile) GetUserWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	currencies, err := p.userApp.GetUserWatchlist(id)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(currencies)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
}

func (p *Profile) createOldPassError(w http.ResponseWriter) {
	var errs entity.ErrorJSON
	errs.Password = append(errs.Password, "введен неверно старый пароль")
	errs.NotEmpty = true

	res, err := json.Marshal(errs)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
	return
}
