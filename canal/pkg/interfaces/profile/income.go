package profile

import (
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	json "github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (p *Profile) GetIncome(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetIncome")

	vars := mux.Vars(r)
	var in = entity.Income{Id: r.Context().Value(entity.UserIdKey).(int64), Period: vars["period"]}

	desc, result, err := p.paymentLogic.GetIncome(r.Context(), in)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	change, err := result.MarshalJSON()
	if err != nil {
		desc = entity.Description{Function: "GetIncome", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(change); err != nil {
		p.logger.Error(entity.Description{Function: "GetIncome", Action: "Write"}, err)
	}
}

func (p *Profile) GetAllIncomePerDay(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetAllIncomePerDay")

	desc, wallets, err := p.paymentLogic.ReceiveWallets(r.Context(), r.Context().Value(entity.UserIdKey).(int64))
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(wallets)
	if err != nil {
		desc = entity.Description{Function: "GetWallets", Action: "Marshal", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err = w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetTransactions", Action: "Write"}, err)
	}
}

func (p *Profile) UpdatePortfolios() {
	task := gocron.NewScheduler(time.UTC)
	defer task.Stop()

	if _, err := task.Every(1).Day().At("00:00").StartImmediately().Do(p.paymentLogic.WritePortfolios); err != nil {
		logrus.WithFields(logrus.Fields{"function": "UpdatePortfolios", "action": "WritePortfolios"}).Error(err)
		return
	}

	<-task.StartAsync()
}
