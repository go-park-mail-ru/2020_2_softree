package profile

import (
	"context"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
	"time"
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

func (p *Profile) writePortfolios() {
	ctx := context.Background()

	userNum, err := p.profile.GetUsers(ctx, &profile.Empty{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "WritePortfolio",
			"action":   "GetUsers",
		}).Error(err)

		return
	}

	for i := int64(0); i < userNum.Num; i++ {
		portfolioValue, err := p.transformActualUserWallets(ctx, i)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "WritePortfolio",
				"action":   "transformActualUserWallets",
				"user_id":  i,
			}).Error(err)

			return
		}

		value, _ := portfolioValue.Float64()
		_, err = p.profile.PutPortfolio(ctx, &profile.PortfolioValue{Id: i, Value: value})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "WritePortfolio",
				"action":   "PutPortfolio",
				"user_id":  i,
				"value":    value,
			}).Error(err)

			return
		}
	}
}

func (p *Profile) UpdatePortfolios() {
	task := gocron.NewScheduler(time.UTC)
	defer task.Stop()

	if _, err := task.Every(1).
		Day().At("00:00").StartImmediately().Do(p.writePortfolios); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "UpdatePortfolios",
		}).Error(err)
		return
	}

	<-task.StartAsync()
}
