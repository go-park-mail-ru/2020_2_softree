package application

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"server/canal/pkg/domain/entity"
	currency "server/currency/pkg/currency/gen"
	"server/currency/pkg/infrastructure/persistence"
)

type CurrencyApp struct {
	currency currency.CurrencyServiceClient
}

func (currencyApp *CurrencyApp) GetAllLatestCurrencies(ctx context.Context) (entity.Description, entity.Currencies, error) {
	out, err := currencyApp.currency.GetAllLatestRates(ctx, &currency.Empty{})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "GetAllLatestCurrencies",
			Action:   "GetAllLatestRates",
		}, entity.Currencies{}, err
	}

	return entity.Description{}, entity.ConvertFromCurrencyCurrencies(out), nil
}

func (currencyApp *CurrencyApp) GetURLCurrencies(r *http.Request) (entity.Description, entity.Currencies, error) {
	vars := mux.Vars(r)
	title := vars["title"]
	if !validateTitle(title) {
		return entity.Description{
			Status:   http.StatusBadRequest,
			Function: "GetURLCurrencies",
			Action:   "validateTitle",
		}, entity.Currencies{}, errors.New("validateTitle from GetURLCurrencies")
	}

	out, err := currencyApp.currency.GetAllRatesByTitle(r.Context(), &currency.CurrencyTitle{Title: title})
	if err != nil {
		return entity.Description{
			Status:   http.StatusInternalServerError,
			Function: "GetAllLatestCurrencies",
			Action:   "GetAllRatesByTitle",
		}, entity.Currencies{}, err
	}

	return entity.Description{}, entity.ConvertFromCurrencyCurrencies(out), nil
}

func (currencyApp *CurrencyApp) GetMarkets() (entity.Description, entity.Markets, error) {
	return entity.Description{}, entity.CreateMarkets(), nil
}

func validateTitle(title string) bool {
	lenOfCurrency := 3
	if len(title) != lenOfCurrency {
		return false
	}

	for _, rate := range persistence.ListOfCurrencies {
		if rate == title {
			return true
		}
	}

	return false
}
