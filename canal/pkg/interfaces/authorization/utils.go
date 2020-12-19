package authorization

import (
	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

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

func (a *Authentication) handleErrorJSON(desc entity.Description, w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(desc.ErrorJSON)
	if err != nil {
		a.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(desc.Status)

	metric.RecordHitMetric(desc.Status, r.URL.Path)

	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "handleErrorJSON",
			"action":   "Write",
		}).Error(err)
	}
}
