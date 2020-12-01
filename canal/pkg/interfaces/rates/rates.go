package rates

import (
	"github.com/prometheus/client_golang/prometheus"
	currency "server/currency/pkg/currency/gen"
)

type Rates struct {
	currencyService currency.CurrencyServiceClient
	Hits            prometheus.CounterVec
}

func NewRates(currencyService currency.CurrencyServiceClient) *Rates {
	return &Rates{
		currencyService: currencyService,
		Hits:            *prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hits"}, []string{"status"}),
	}
}
