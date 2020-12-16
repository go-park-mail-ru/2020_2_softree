package profile

import (
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"server/canal/pkg/domain/entity"
	"time"
)

func (p *Profile) GetIncome(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var in = entity.Income{Id: r.Context().Value(entity.UserIdKey).(int64), Period: vars["period"]}

	desc, result, err := p.paymentLogic.GetIncome(r.Context(), in)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	change, err := result.MarshalJSON()
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetIncome", Action: "Marshal", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	p.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(change); err != nil {
		p.logger.Error(entity.Description{Function: "GetIncome", Action: "Write"}, err)
	}
}

func (p *Profile) UpdatePortfolios() {
	task := gocron.NewScheduler(time.UTC)
	defer task.Stop()

	if _, err := task.Every(1).
		Day().At("00:00").StartImmediately().Do(p.paymentLogic.WritePortfolios); err != nil {
		logrus.WithFields(logrus.Fields{"function": "UpdatePortfolios"}).Error(err)
		return
	}

	<-task.StartAsync()
}
