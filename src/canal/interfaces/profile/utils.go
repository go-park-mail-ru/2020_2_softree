package profile

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
	"server/src/canal/domain/entity"
	profile "server/src/profile/profile/gen"

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

func (p *Profile) checkWalletFrom(ctx context.Context, wallet *profile.ConcreteWallet) (bool, int) {
	var exist *profile.Check
	var err error
	if exist, err = p.userApp.CheckWallet(ctx, wallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletFrom",
			"action":   "CheckWallet",
		}).Error(err)
		return false, http.StatusInternalServerError
	}

	if !exist.Existence {
		return false, http.StatusBadRequest
	}

	return true, 0
}

func (p *Profile) checkWalletTo(ctx context.Context, wallet *profile.ConcreteWallet) (bool, int) {
	var exist *profile.Check
	var err error
	if exist, err = p.userApp.CheckWallet(ctx, wallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletTo",
			"action":   "CheckWallet",
		}).Error(err)
		return false, http.StatusInternalServerError
	}

	if !exist.Existence {
		if _, err = p.userApp.CreateWallet(ctx, wallet); err != nil {
			logrus.WithFields(logrus.Fields{
				"status":   http.StatusInternalServerError,
				"function": "checkWallets",
				"action":   "CreateWallet",
			}).Error(err)
			return false, http.StatusInternalServerError
		}
	}

	return true, 0
}

func (p *Profile) getCurrencyDiv(
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

const notEnoughPayment = 1

func (p *Profile) getPay(ctx context.Context, userWallet *profile.ConcreteWallet, needToPay decimal.Decimal) int {
	var wallet *profile.Wallet
	var err error
	if wallet, err = p.userApp.GetWallet(ctx, userWallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "getPay",
			"transactionFrom": userWallet.Title,
			"action":          "GetWallet",
		}).Error(err)
		return http.StatusInternalServerError
	}

	if needToPay.GreaterThan(decimal.NewFromFloat(wallet.Value)) {
		return notEnoughPayment
	}

	return 0
}

func (p *Profile) validate(action string, user *profile.User) bool {
	switch action {
	case "Avatar":
		if govalidator.IsNull(user.Avatar) {
			logrus.WithFields(logrus.Fields{
				"status":   http.StatusBadRequest,
				"function": "UpdateUserAvatar",
				"action":   "validation",
			}).Error("No user avatar from json")
			return false
		}
	case "Passwords":
		if govalidator.IsNull(user.OldPassword) || govalidator.IsNull(user.NewPassword) {
			logrus.WithFields(logrus.Fields{
				"status":      http.StatusBadRequest,
				"function":    "UpdateUserPassword",
				"oldPassword": user.OldPassword,
				"newPassword": user.NewPassword,
			}).Error("No user passwords from json")
			return false
		}
	}

	return true
}

func (p *Profile) ValidateUpdate(u *profile.User) (errors entity.ErrorJSON) {
	if govalidator.HasWhitespace(u.NewPassword) {
		errors.Password = append(errors.Email, "Некорректный новый пароль")
		errors.NotEmpty = true
	}

	if govalidator.HasWhitespace(u.OldPassword) {
		errors.Password = append(errors.Email, "Некорректный старый пароль")
		errors.NotEmpty = true
	}

	return errors
}
