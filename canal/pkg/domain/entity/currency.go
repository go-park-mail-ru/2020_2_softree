package entity

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/shopspring/decimal"
	"log"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"
	"time"
)

//easyjson:json
type (
	Currency struct {
		Base      string          `json:"base"`
		Title     string          `json:"title"`
		Value     decimal.Decimal `json:"value"`
		UpdatedAt time.Time       `json:"updated_at"`
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
		updated, err := ptypes.Timestamp(currency.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}

		currenciesEntity = append(currenciesEntity, Currency{
			Base:  currency.Base,
			Title: currency.Title,
			Value: decimal.NewFromFloat(currency.Value),
			UpdatedAt: updated,
		})
	}

	return currenciesEntity
}
