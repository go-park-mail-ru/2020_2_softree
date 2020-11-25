package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (p *Profile) GetTransactions(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	history, err := p.profile.GetAllPaymentHistory(r.Context(), &profile.UserID{Id: id})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "GetAllPaymentHistory",
			"userID":   id,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(history)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "Marshal",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetTransactions",
			"action":   "Write",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Profile) SetTransaction(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int64)

	var transaction profile.PaymentHistory
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "Decode",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if exist, code := p.checkWalletFrom(r.Context(), &profile.ConcreteWallet{Id: id, Title: transaction.From}); !exist {
		w.WriteHeader(code)
		return
	}

	var div decimal.Decimal
	// #TODO
	if err, div = p.getCurrencyDiv(); err != nil {
		return
	}

	needToPay := div.Mul(decimal.NewFromFloat(transaction.Amount))
	if code := p.getPay(r.Context(), &profile.ConcreteWallet{Id: id, Title: transaction.To}, needToPay); code != 0 {
		if code == notEnoughPayment {
			p.createErrorJSON(errors.New("not enough payment"))
			return
		}
		return
	}

	if exist, code := p.checkWalletTo(r.Context(), &profile.ConcreteWallet{Id: id, Title: transaction.To}); !exist {
		w.WriteHeader(code)
		return
	}

	needToPay = needToPay.Mul(decimal.New(-1, 0))
	money, _ := needToPay.Float64()
	toSetWallet := profile.ToSetWallet{Id: id, NewWallet: &profile.Wallet{Title: transaction.From, Value: money}}
	if _, err = p.profile.UpdateWallet(r.Context(), &toSetWallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "UpdateWallet",
			"title":    transaction.From,
			"value":    money,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	toSetWallet = profile.ToSetWallet{Id: id, NewWallet: &profile.Wallet{Title: transaction.To, Value: transaction.Amount}}
	if _, err = p.profile.UpdateWallet(r.Context(), &toSetWallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "SetTransactions",
			"action":   "UpdateWallet",
			"title":    transaction.To,
			"value":    transaction.Amount,
		}).Error(err)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
