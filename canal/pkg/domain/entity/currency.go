package entity

import profile "server/profile/pkg/profile/gen"

type Currency struct {
	Base  string `json:"base"`
	Title string `json:"title"`
}

func ConvertToSlice(currenciesProfile *profile.Currencies) []Currency {
	currenciesEntity := make([]Currency, 0, len(currenciesProfile.Currencies))
	for _, currency := range currenciesProfile.Currencies {
		currenciesEntity = append(currenciesEntity, Currency{Base: currency.Base, Title: currency.Title})
	}

	return currenciesEntity
}
