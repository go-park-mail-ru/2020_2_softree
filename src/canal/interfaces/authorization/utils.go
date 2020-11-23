package authorization

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"server/src/canal/domain/entity"
	"server/src/canal/infrastructure/security"
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
			"status": http.StatusInternalServerError}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": http.StatusInternalServerError}).Error(err)
	}
}

func CreateCookie() (http.Cookie, error) {
	hash, err := security.MakeShieldedCookie()
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   viper.GetString("server.domain"),
		Secure:   viper.GetBool("server.secure"),
		HttpOnly: true,
		Path:     "/",
	}, nil
}
