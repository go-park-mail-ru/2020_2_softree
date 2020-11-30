package authorization

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func (a *Authentication) checkGetUserByLoginErrors(err error) (errs entity.ErrorJSON) {
	if err == nil {
		return
	}

	if err.Error() == "wrong password" {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Неправильный email или пароль")
	}

	return
}

func (a *Authentication) createServerError(errors *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "createServerError",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createServerError",
			"action":   "Write",
		}).Error(err)
	}
}

func (a *Authentication) validate(user *profile.User) (errs entity.ErrorJSON) {
	if !govalidator.IsEmail(user.Email) {
		errs.Email = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if user.Password == "" {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if govalidator.IsNull(user.Password) {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	if govalidator.HasWhitespace(user.Password) {
		errs.Password = append(errs.Email, "Некорректный email или пароль")
		errs.NotEmpty = true
	}

	return
}

func (a *Authentication) recordHitMetric(code int, url string) {
	a.hits.WithLabelValues(strconv.Itoa(code), url)
}

func CreateCookie() http.Cookie {
	return http.Cookie{
		Name:     "session_id",
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   viper.GetString("server.domain"),
		Secure:   viper.GetBool("server.secure"),
		HttpOnly: true,
		Path:     "/",
	}
}
