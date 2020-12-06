package profile

import (
	"context"
	json "github.com/mailru/easyjson"
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
		code := http.StatusInternalServerError
		desc := entity.Description{Function: "createServerError", Action: "Marshal", Err: err, Status: code}
		p.logger.Error(desc)
		w.WriteHeader(http.StatusInternalServerError)

		p.recordHitMetric(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	p.recordHitMetric(http.StatusBadRequest)
	if _, err := w.Write(res); err != nil {
		p.logger.Error(entity.Description{Function: "UpdateUserPassword", Action: "Write", Err: err})
	}
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
