package profile

import (
	"context"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"
	"server/canal/pkg/domain/entity"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"
	"strconv"

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

		p.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	p.recordHitMetric(http.StatusBadRequest)

	if _, err := w.Write(res); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "createServerError",
			"action":   "Write",
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

		p.recordHitMetric(http.StatusInternalServerError)

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

		p.recordHitMetric(http.StatusInternalServerError)

		return false, http.StatusInternalServerError
	}

	if !exist.Existence {
		if _, err = p.profile.CreateWallet(ctx, wallet); err != nil {
			logrus.WithFields(logrus.Fields{
				"status":   http.StatusInternalServerError,
				"function": "checkWallets",
				"action":   "CreateWallet",
			}).Error(err)

			p.recordHitMetric(http.StatusInternalServerError)

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
			"status":   http.StatusInternalServerError,
			"function": "checkWalletPayment",
			"Base":     transaction.Base,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

		return err, http.StatusInternalServerError, decimal.Decimal{}
	}

	var currencyCurr *currency.Currency
	if currencyCurr, err = p.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: transaction.Currency}); err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "checkWalletPayment",
			"Currency": transaction.Currency,
		}).Error(err)

		p.recordHitMetric(http.StatusInternalServerError)

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

		p.recordHitMetric(http.StatusInternalServerError)

		return http.StatusInternalServerError
	}

	if needToPay.GreaterThan(decimal.NewFromFloat(wallet.Value)) {
		return notEnoughPayment
	}

	return 0
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

func (p *Profile) recordHitMetric(code int) {
	p.Hits.WithLabelValues(strconv.Itoa(code)).Inc()
}

func (p *Profile) transformActualUserWallets(ctx context.Context, id int64) (decimal.Decimal, error) {
	wallets, err := p.profile.GetWallets(ctx, &profile.UserID{Id: id})
	if err != nil {
		return decimal.Decimal{}, err
	}

	var cash decimal.Decimal
	for _, wallet := range wallets.Wallets {
		curr, err := p.rates.GetLastRate(ctx, &currency.CurrencyTitle{Title: wallet.Title})
		if err != nil {
			return decimal.Decimal{}, err
		}
		cash = cash.Add(decimal.NewFromFloat(wallet.Value / curr.Value))
	}

	return cash, nil
}
