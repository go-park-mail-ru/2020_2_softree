package profile

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/infrastructure/logger"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	history, err := p.userApp.GetAllPaymentHistory(id)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"userID":   id,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var transaction entity.PaymentHistory
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if !p.checkWalletFrom(w, id, transaction) {
		return
	}

	var div decimal.Decimal
	if err, div = p.checkWalletPayment(w, transaction); err != nil {
		if err.Error() == "not enough payment" {
			p.createErrorJSON(err)
			return
		}
		return
	}

	var val decimal.Decimal
	if err, val = p.getPay(w, id, transaction, div); err != nil {
		if err.Error() == "not enough payment" {
			p.createErrorJSON(err)
			return
		}
		return
	}

	if !p.checkWalletTo(w, id, transaction) {
		return
	}

	val = val.Mul(decimal.New(-1, 0))
	if err = p.userApp.UpdateWallet(id, entity.Wallet{Title: transaction.From, Value: val}); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = p.userApp.UpdateWallet(id, entity.Wallet{Title: transaction.To, Value: transaction.Amount}); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transaction.Value = div
	if err = p.userApp.AddToPaymentHistory(id, transaction); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p *Profile) GetHistoryInInterval(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(uint64)

	var interval entity.Interval
	err := json.NewDecoder(r.Body).Decode(&interval)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetHistoryInInterval",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	history, err := p.userApp.GetIntervalPaymentHistory(id, interval)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetHistoryInInterval",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history)
	if err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetHistoryInInterval",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logger.GlobalLogger.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetHistoryInInterval",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
