package profile

import (
	"errors"
	json "github.com/mailru/easyjson"
	"io/ioutil"
	"net/http"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
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
	id := r.Context().Value(entity.UserIdKey).(int64)

	var transaction profile.PaymentHistory
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "ReadAll",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &transaction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "Unmarshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	var div decimal.Decimal
	var code int
	if err, code, div = p.getCurrencyDiv(r.Context(), &transaction); err != nil {
		w.WriteHeader(code)
		p.recordHitMetric(code)
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
		w.WriteHeader(code)
		p.recordHitMetric(code)
		return
	}

	if code := p.getPay(r.Context(), &profile.ConcreteWallet{Id: id, Title: titleToCheckPayment}, checkingPayment); code != 0 {
		if code == notEnoughPayment {
			errs := p.createErrorJSON(errors.New("not enough payment"))
			p.createServerError(&errs, w)
			return
		}
		w.WriteHeader(code)

		p.recordHitMetric(code)
		return
	}

	if exist, code := p.checkWalletBuy(r.Context(), &profile.ConcreteWallet{Id: id, Title: putTitle}); !exist {
		w.WriteHeader(code)
		p.recordHitMetric(code)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	p.recordHitMetric(http.StatusCreated)
}
