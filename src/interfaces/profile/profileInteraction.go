package profile

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
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

func (p *Profile) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !govalidator.IsBase64(user.Avatar) {
		p.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.sanitizer.SanitizeBytes([]byte(user.Avatar))
	if err = p.userApp.UpdateUserAvatar(id, user); err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user, err = p.userApp.GetUserById(id); err != nil {
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
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
}

func (p *Profile) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if govalidator.IsNull(user.OldPassword) || govalidator.IsNull(user.NewPassword) {
		p.log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.sanitizer.Sanitize(user.OldPassword)
	p.sanitizer.Sanitize(user.NewPassword)

	errs := user.ValidateUpdate()
	if errs.NotEmpty {
		p.createServerError(&errs, w)
		return
	}

	var check bool
	if check, err = p.userApp.CheckPassword(id, user.OldPassword); err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !check {
		errs := p.createErrorJSON(errors.New("wrong old password"))
		p.createServerError(&errs, w)
		return
	}

	if err = p.userApp.UpdateUserPassword(id, user); err != nil {
		p.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user, err = p.userApp.GetUserById(id); err != nil {
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
	if _, err := w.Write(res); err != nil {
		p.log.Print(err)
	}
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
