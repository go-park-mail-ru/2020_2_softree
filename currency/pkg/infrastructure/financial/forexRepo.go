package financial

import (
	"server/currency/pkg/infrastructure/persistence"
)

type ForexRepo struct {
	forex map[string]float64
	base  string
}

func convertToForexRepo(rates map[string]float64) *ForexRepo {
	finance := &ForexRepo{
		forex: make(map[string]float64, persistence.LenListOfCurrencies),
		base:  "USD",
	}
	for key, val := range rates {
		finance.forex[key] = val
	}

	return finance
}

func (fr *ForexRepo) GetBase() string {
	return fr.base
}

func (fr *ForexRepo) GetQuote() map[string]float64 {
	return fr.forex
}
