package financial

import (
	"fmt"
	"server/currency/pkg/infrastructure/persistence"
	"strconv"
)

type ForexRepo struct {
	forex map[string]interface{}
	base  string
}

func convertToForexRepo(rates map[string]string) *ForexRepo {
	finance := &ForexRepo{
		forex: make(map[string]interface{}, persistence.LenListOfCurrencies),
		base:  "USD",
	}
	var err error
	for key, val := range rates {
		finance.forex[key], err = strconv.ParseFloat(val, 3)
		if err != nil {
			fmt.Println(err)
		}
	}

	return finance
}

func (fr *ForexRepo) GetBase() string {
	return fr.base
}

func (fr *ForexRepo) GetQuote() map[string]interface{} {
	return fr.forex
}
