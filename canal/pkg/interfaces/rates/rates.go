package rates

import (
	"github.com/prometheus/client_golang/prometheus"
	"server/canal/pkg/infrastructure/metric"
	currency "server/currency/pkg/currency/gen"
)

type Rates struct {
	currencyService currency.CurrencyServiceClient
	Hits            prometheus.CounterVec
}

func NewRates(currencyService currency.CurrencyServiceClient) *Rates {
	return &Rates{
		currencyService: currencyService,
		Hits:            *metric.Metric,
	}
}
