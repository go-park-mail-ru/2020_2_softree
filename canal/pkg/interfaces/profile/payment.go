package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, payments := p.paymentLogic.ReceiveTransactions(r.Context(), id)
	if desc.Err != nil {
		p.logger.Error(desc)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(payments)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetTransactions", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetTransactions", Action: "Write", Err: err})
	}
}

func (p *Profile) SetTransaction(w http.ResponseWriter, r *http.Request) {
	transaction, desc := entity.GetTransactionFromBody(r.Body)
	if desc.Err != nil {
		desc.Function = "SetTransaction"
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	transaction.UserId = r.Context().Value(entity.UserIdKey).(int64)

	if desc := p.paymentLogic.SetPayment(r.Context(), transaction); desc.Err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "SetTransaction", Action: "SetPayment", Err: desc.Err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	p.recordHitMetric(http.StatusCreated)
}
