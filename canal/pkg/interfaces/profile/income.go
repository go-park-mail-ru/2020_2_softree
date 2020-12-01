package profile

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
)

func (p *Profile) GetIncome(w http.ResponseWriter, r *http.Request) {
	var incomeParameters profile.IncomeParameters
	err := json.NewDecoder(r.Body).Decode(&incomeParameters.Period)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetIncome",
			"action":   "Decode",
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	incomeParameters.Id = r.Context().Value(entity.UserIdKey).(int64)

	result, err := p.profile.GetIncome(r.Context(), &incomeParameters)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetIncome",
			"action":   "GetIncome",
			"incomeParameters":   &incomeParameters,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	change, err := json.Marshal(result.Change)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetWallets",
			"action":   "Marshal",
			"change":   change,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	p.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(change); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetIncome",
			"action":   "Write",
			"change":   change,
		}).Error(err)
		return
	}
}
