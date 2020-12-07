package profile

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	desc, payments, err := p.paymentLogic.ReceiveTransactions(r.Context(), id)
	if err != nil {
		p.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		p.recordHitMetric(desc.Status)
		return
	}

	res, err := json.Marshal(payments)
	if err != nil {
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "GetTransactions", Action: "Marshal", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	p.recordHitMetric(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "GetTransactions", Action: "Write"}, err)
	}
}

func (p *Profile) SetTransaction(w http.ResponseWriter, r *http.Request) {
	transaction, desc, err := entity.GetTransactionFromBody(r.Body)
	if err != nil {
		desc.Function = "SetTransaction"
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	transaction.UserId = r.Context().Value(entity.UserIdKey).(int64)

	if desc, err = p.paymentLogic.SetTransaction(r.Context(), transaction); err != nil {
		code := http.StatusInternalServerError
		desc = entity.Description{Function: "SetTransaction", Action: "SetPayment", Status: code}
		p.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	p.recordHitMetric(http.StatusCreated)
}
