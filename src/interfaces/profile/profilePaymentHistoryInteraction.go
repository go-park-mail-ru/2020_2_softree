package profile

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	history, err := p.userApp.GetAllPaymentHistory(id)
	if err != nil {
		p.log.Info("user id: ", id, ", func: GetAllPaymentHistory, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history)
	if err != nil {
		p.log.Info("func: GetHistory, with error while marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.log.Info("func: GetAllPaymentHistory, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var transaction entity.PaymentHistory
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		p.log.Info("func: SetTransactions, with error while decode json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !p.checkWalletFrom(w, id, transaction) {
		return
	}

	if err = p.checkWalletPayment(w, id, transaction); err != nil {
		if err.Error() == "not enough payment" {
			p.createErrorJSON(err)
			return
		}
		return
	}
}

func (p *Profile) GetHistoryInInterval(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var interval entity.Interval
	err := json.NewDecoder(r.Body).Decode(&interval)
	if err != nil {
		p.log.Info("func: GetHistoryInInterval, with error while decode json: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	history, err := p.userApp.GetIntervalPaymentHistory(id, interval)
	if err != nil {
		p.log.Info("func: GetHistoryInInterval, with GetIntervalPaymentHistory error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history)
	if err != nil {
		p.log.Info("func: GetHistoryInInterval, with error while marshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		p.log.Info("func: GetHistoryInInterval, with error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
