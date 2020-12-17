package entity

import (
	profile "server/profile/pkg/profile/gen"
	currency "server/currency/pkg/currency/gen"
)

//easyjson:json
type (
	Currency struct {
		Base  string `json:"base"`
		Title string `json:"title"`
	}

	Currencies []Currency
)

func ConvertFromProfileCurrencies(currenciesProfile *profile.Currencies) Currencies {
	currenciesEntity := make(Currencies, 0, len(currenciesProfile.Currencies))
	for _, currency := range currenciesProfile.Currencies {
		currenciesEntity = append(currenciesEntity, Currency{
			Base:  currency.Base,
			Title: currency.Title,
		})
	}

	return currenciesEntity
}

func ConvertFromCurrencyCurrencies(currenciesCurrency *currency.Currencies) Currencies {
	currenciesEntity := make(Currencies, 0, len(currenciesCurrency.Rates))
	for _, currency := range currenciesCurrency.Rates {
		currenciesEntity = append(currenciesEntity, Currency{
			Base:  currency.Base,
			Title: currency.Title,
		})
	}

	return currenciesEntity
}
