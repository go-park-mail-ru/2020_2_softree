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


func (a *Authentication) createServerError(errors *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errors)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "createServerError",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		a.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	a.recordHitMetric(http.StatusBadRequest)

	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createServerError",
			"action":   "Write",
		}).Error(err)
	}
}

func (a *Authentication) recordHitMetric(code int) {
	a.Hits.WithLabelValues(strconv.Itoa(code)).Inc()
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
