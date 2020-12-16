package entity

import profile "server/profile/pkg/profile/gen"

type Currency struct {
	Base  string `json:"base"`
	Title string `json:"title"`
}

//easyjson:json
type Currencies []Currency

func ConvertToSlice(currenciesProfile *profile.Currencies) Currencies {
	currenciesEntity := make(Currencies, 0, len(currenciesProfile.Currencies))
	for _, currency := range currenciesProfile.Currencies {
		currenciesEntity = append(currenciesEntity, Currency{Base: currency.Base, Title: currency.Title})
	}

	return currenciesEntity
}
