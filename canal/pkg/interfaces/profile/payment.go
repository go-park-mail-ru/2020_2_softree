package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	history, err := p.profile.GetAllPaymentHistory(r.Context(), &profile.UserID{Id: id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "GetAllPaymentHistory",
			"userID":   id,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history.History)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "Marshal",
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "Write",
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(entity.UserIdKey).(int64)

	var transaction profile.PaymentHistory
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "Decode",
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var div decimal.Decimal
	var code int
	if err, code, div = p.getCurrencyDiv(r.Context(), &transaction); err != nil {
		p.recordHitMetric(code)
		w.WriteHeader(code)
		return
	}

	divMulAmount := div.Mul(decimal.NewFromFloat(transaction.Amount))

	titleToCheckPayment := transaction.Currency
	checkingPayment := decimal.NewFromFloat(transaction.Amount)

	removedMoney := -transaction.Amount
	removedTitle := transaction.Currency

	putMoney, _ := divMulAmount.Float64()
	putTitle := transaction.Base
	if transaction.Sell == "false" {
		titleToCheckPayment = transaction.Base
		checkingPayment = divMulAmount

		removedMoney, _ = checkingPayment.Float64()
		removedMoney *= -1
		removedTitle = transaction.Base

		putMoney = transaction.Amount
		putTitle = transaction.Currency
	}

	if exist, code := p.checkWalletSell(r.Context(), &profile.ConcreteWallet{Id: id, Title: titleToCheckPayment}); !exist {
		p.recordHitMetric(code)
		w.WriteHeader(code)
		return
	}

	if code := p.getPay(r.Context(), &profile.ConcreteWallet{Id: id, Title: titleToCheckPayment}, checkingPayment); code != 0 {
		if code == notEnoughPayment {
			errs := p.createErrorJSON(errors.New("not enough payment"))
			p.createServerError(&errs, w)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(code)
		return
	}

	if exist, code := p.checkWalletBuy(r.Context(), &profile.ConcreteWallet{Id: id, Title: putTitle}); !exist {
		p.recordHitMetric(code)
		w.WriteHeader(code)
		return
	}

	toSetWallet := profile.ToSetWallet{Id: id, NewWallet: &profile.Wallet{Title: removedTitle, Value: removedMoney}}
	if _, err = p.profile.UpdateWallet(r.Context(), &toSetWallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "UpdateWallet",
			"title":    removedTitle,
			"value":    removedMoney,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	toSetWallet = profile.ToSetWallet{Id: id, NewWallet: &profile.Wallet{Title: putTitle, Value: putMoney}}
	if _, err = p.profile.UpdateWallet(r.Context(), &toSetWallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "UpdateWallet",
			"title":    putTitle,
			"value":    putMoney,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	transaction.Value, _ = div.Float64()
	if _, err = p.profile.AddToPaymentHistory(r.Context(), &profile.AddToHistory{Id: id, Transaction: &transaction}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":      http.StatusInternalServerError,
			"function":    "SetTransactions",
			"action":      "AddToPaymentHistory",
			"id":          id,
			"transaction": &transaction,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p.recordHitMetric(http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
}
