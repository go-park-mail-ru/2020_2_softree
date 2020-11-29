package profile

import (
	"encoding/json"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"

	"github.com/sirupsen/logrus"
)

func (p *Profile) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	var user profile.User
	var err error
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserAvatar",
			"action":   "Decode",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)
	defer r.Body.Close()

	if !p.validate("Avatar", &user) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = p.profile.UpdateUserAvatar(r.Context(), &profile.UpdateFields{Id: user.Id, User: &user}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserAvatar",
			"action":   "UpdateUserAvatar",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	public, err := p.profile.GetUserById(r.Context(), &profile.UserID{Id: user.Id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserAvatar",
			"action":   "GetUserById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.sanitizer.SanitizeBytes([]byte(public.Avatar))
	res, err := json.Marshal(public)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserAvatar",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateUserAvatar",
		}).Error(err)
	}
}

func (p *Profile) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var user profile.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "Decode",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Id = r.Context().Value(entity.UserIdKey).(int64)

	if !p.validate("Passwords", &user) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.sanitizer.Sanitize(user.OldPassword)
	p.sanitizer.Sanitize(user.NewPassword)

	if errs := p.validateUpdate(&user); errs.NotEmpty {
		p.createServerError(&errs, w)
		return
	}

	var check *profile.Check
	if check, err = p.profile.CheckPassword(r.Context(), &user); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "CheckPassword",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check.Existence {
		p.createOldPassError(w)
		return
	}

	if user.NewPassword, err = p.security.MakeShieldedPassword(user.NewPassword); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "MakeShieldedPassword",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = p.profile.UpdateUserPassword(r.Context(), &profile.UpdateFields{Id: user.Id, User: &user}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "UpdateUserPassword",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var public *profile.PublicUser
	if public, err = p.profile.GetUserById(r.Context(), &profile.UserID{Id: user.Id}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "GetUserById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "UpdateUserPassword",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdateUserPassword",
			"action":   "Write",
		}).Error(err)
	}
}

func (p *Profile) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	var err error
	var public *profile.PublicUser
	if public, err = p.profile.GetUserById(r.Context(), &profile.UserID{Id: id}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetUser",
			"action":   "GetUserById",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(public)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetUser",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetUser",
			"action":   "Write",
		}).Error(err)
	}
}

func (p *Profile) GetUserWatchlist(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	currencies, err := p.profile.GetUserWatchlist(r.Context(), &profile.UserID{Id: id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusBadRequest,
			"function": "GetUserWatchlist",
			"action":   "GetUserWatchlist",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := json.Marshal(currencies)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetUserWatchlist",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetUserWatchlist",
			"action":   "Write",
		}).Error(err)
	}
}

func (p *Profile) createOldPassError(w http.ResponseWriter) {
	var errs entity.ErrorJSON
	errs.Password = append(errs.Password, "введен неверно старый пароль")
	errs.NotEmpty = true

	res, err := json.Marshal(errs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "createOldPassError",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createOldPassError",
			"action":   "Write",
		}).Error(err)
	}
}
