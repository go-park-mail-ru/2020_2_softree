package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/src/domain/entity"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

func (p *Profile) createErrorJSON(err error) (errs entity.ErrorJSON) {
	if err.Error() == "wrong old password" {
		errs.NotEmpty = true
		errs.NonFieldError = append(errs.NonFieldError, "Введен неверно старый пароль")
	}

	if err.Error() == "not enough payment" {
		errs.NotEmpty = true
		errs.NonFieldError = append(
			errs.NonFieldError,
			"В вашем кошельке недостаточно средств для совершения данной транзакции",
		)
	}

	return
}

func (p *Profile) createServerError(errs *entity.ErrorJSON, w http.ResponseWriter) {
	res, err := json.Marshal(errs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "createServerError",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createServerError",
		}).Error(err)
	}
}

func (p *Profile) checkWalletFrom(w http.ResponseWriter, id uint64, transaction entity.PaymentHistory) bool {
	var exist bool
	var err error
	if exist, err = p.userApp.CheckWallet(id, transaction.From); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletFrom",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !exist {
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}

func (p *Profile) checkWalletTo(w http.ResponseWriter, id uint64, transaction entity.PaymentHistory) bool {
	var exist bool
	var err error
	if exist, err = p.userApp.CheckWallet(id, transaction.To); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletTo",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if !exist {
		if err = p.userApp.CreateWallet(id, transaction.To); err != nil {
			logrus.WithFields(logrus.Fields{
				"status":   http.StatusInternalServerError,
				"function": "checkWallets",
			}).Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return false
		}
	}

	return true
}

func (p *Profile) checkWalletPayment(
	w http.ResponseWriter, transaction entity.PaymentHistory) (error, decimal.Decimal) {
	var currencyFrom entity.Currency
	var err error
	if currencyFrom, err = p.rateApp.GetLastRate(transaction.From); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "checkWalletPayment",
			"transactionFrom": transaction.From,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, decimal.Decimal{}
	}

	var currencyTo entity.Currency
	if currencyTo, err = p.rateApp.GetLastRate(transaction.To); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "checkWalletPayment",
			"transactionFrom": transaction.To,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, decimal.Decimal{}
	}

	div := currencyFrom.Value.Div(currencyTo.Value)
	return nil, div
}

func (p *Profile) getPay(
	w http.ResponseWriter, id uint64, transaction entity.PaymentHistory, div decimal.Decimal) (err error, needToPay decimal.Decimal) {
	needToPay = div.Mul(transaction.Amount)

	var wallet entity.Wallet
	if wallet, err = p.userApp.GetWallet(id, transaction.From); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "getPay",
			"transactionFrom": transaction.From,
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return err, decimal.Decimal{}
	}

	if needToPay.GreaterThan(wallet.Value) {
		return errors.New("not enough payment"), decimal.Decimal{}
	}

	return nil, needToPay
}
