package rates

import (
	currency "server/currency/pkg/currency/gen"
)

type Rates struct {
	currencyService currency.CurrencyServiceClient
}

func NewRates(currencyService currency.CurrencyServiceClient) *Rates {
	return &Rates{currencyService}
}