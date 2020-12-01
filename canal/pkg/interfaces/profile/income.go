package profile

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
)

func (p *Profile) GetIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]
	var incomeParameters = profile.IncomeParameters{Id: r.Context().Value(entity.UserIdKey).(int64), Period: period}

	result, err := p.profile.GetIncome(r.Context(), &incomeParameters)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":           http.StatusInternalServerError,
			"function":         "GetIncome",
			"action":           "GetIncome",
			"incomeParameters": &incomeParameters,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	walletUSDCash, err := p.transformActualUserWallets(r.Context(), incomeParameters.Id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":        http.StatusInternalServerError,
			"function":      "GetIncome",
			"action":        "transformActualUserWallets",
			"walletUSDCash": walletUSDCash,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result.Change, _ = walletUSDCash.Sub(decimal.NewFromFloat(-result.Change)).Float64()

	change, err := json.Marshal(result.Change)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetIncome",
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
