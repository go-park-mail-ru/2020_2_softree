package profile

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
	"server/canal/pkg/domain/entity"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"

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

func (p *Profile) checkWalletSell(ctx context.Context, wallet *profile.ConcreteWallet) (bool, int) {
	var exist *profile.Check
	var err error
	if exist, err = p.profile.CheckWallet(ctx, wallet); err != nil {
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

func (p *Profile) checkWalletBuy(ctx context.Context, wallet *profile.ConcreteWallet) (bool, int) {
	var exist *profile.Check
	var err error
	if exist, err = p.profile.CheckWallet(ctx, wallet); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletTo",
			"action":   "CheckWallet",
		}).Error(err)
		return false, http.StatusInternalServerError
	}

	if !exist.Existence {
		if _, err = p.profile.CreateWallet(ctx, wallet); err != nil {
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
	ctx context.Context, transaction *profile.PaymentHistory) (error, int, decimal.Decimal) {
	var currencyBase *currency.Currency
	var err error
	if currencyBase, err = p.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Base}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "checkWalletPayment",
			"Base": transaction.Base,
		}).Error(err)
		return err, http.StatusInternalServerError, decimal.Decimal{}
	}

	var currencyCurr *currency.Currency
	if currencyCurr, err = p.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Currency}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":          http.StatusInternalServerError,
			"function":        "checkWalletPayment",
			"Currency": transaction.Currency,
		}).Error(err)
		return err, http.StatusInternalServerError, decimal.Decimal{}
	}

	div := decimal.NewFromFloat(currencyBase.Value).Div(decimal.NewFromFloat(currencyCurr.Value))
	return nil, 0, div
}

const notEnoughPayment = 1

func (p *Profile) getPay(ctx context.Context, userWallet *profile.ConcreteWallet, needToPay decimal.Decimal) int {
	var wallet *profile.Wallet
	var err error
	if wallet, err = p.profile.GetWallet(ctx, userWallet); err != nil {
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

func (p *Profile) validateUpdate(u *profile.User) (errors entity.ErrorJSON) {
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
