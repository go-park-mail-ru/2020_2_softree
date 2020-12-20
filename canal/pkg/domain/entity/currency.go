package entity

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shopspring/decimal"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"
)

//easyjson:json
type (
	Currency struct {
		Base      string               `json:"base,omitempty"`
		Title     string               `json:"title"`
		Value     decimal.Decimal      `json:"value"`
		UpdatedAt *timestamp.Timestamp `json:"updated_at,omitempty"`
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
		if currency.Base == "" {
			currenciesEntity = append(currenciesEntity, Currency{
				Title:     currency.Title,
				Value:     decimal.NewFromFloat(currency.Value),
				UpdatedAt: currency.UpdatedAt,
			})
		} else {
			currenciesEntity = append(currenciesEntity, Currency{
				Base:      currency.Base,
				Title:     currency.Title,
				Value:     decimal.NewFromFloat(currency.Value),
				UpdatedAt: currency.UpdatedAt,
			})
		}
	}

	return currenciesEntity
}

func ConvertFromInitialDayCurrencies(currenciesCurrency *currency.InitialDayCurrencies) Currencies {
	currenciesEntity := make(Currencies, 0, len(currenciesCurrency.Currencies))
	for _, currency := range currenciesCurrency.Currencies {
		currenciesEntity = append(currenciesEntity, Currency{
			Title: currency.Title,
			Value: decimal.NewFromFloat(currency.Value),
		})
	}

	return currenciesEntity
}
