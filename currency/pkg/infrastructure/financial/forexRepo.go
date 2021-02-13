package financial

import (
	"server/currency/pkg/infrastructure/persistence"
)

type ForexRepo struct {
	forex map[string]interface{}
	base  string
}

func convertToForexRepo(rates map[string]interface{}) *ForexRepo {
	finance := ForexRepo{
		forex: make(map[string]interface{}, persistence.LenListOfCurrencies),
		base:  "USD",
	}
	for key, val := range rates["quote"].(map[string]interface{}) {
		finance.forex[key] = val
	}

	return &finance
}

func (fr *ForexRepo) GetBase() string {
	return fr.base
}

func (fr *ForexRepo) GetQuote() map[string]interface{} {
	return fr.forex
}
