package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetTransactions")

	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, payments, err := p.paymentLogic.ReceiveTransactions(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	res, err := json.Marshal(payments)
	if err != nil {
		desc := entity.Description{Function: "GetTransactions", Action: "Marshal", Status: http.StatusInternalServerError}
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

func (p *Profile) SetTransaction(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "SetTransaction")

	transaction, desc, err := entity.GetTransactionFromBody(r.Body)
	if err != nil {
		desc.Function = "SetTransaction"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}
	transaction.UserId = r.Context().Value(entity.UserIdKey).(int64)

	if desc, err = p.paymentLogic.SetTransaction(r.Context(), transaction); err != nil {
		desc = entity.Description{Function: "SetTransaction", Action: "SetPayment", Status: http.StatusInternalServerError}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.WriteHeader(http.StatusCreated)
	metric.RecordHitMetric(http.StatusCreated, r.URL.Path)
}
