package authorization

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"

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
