package financial

import "github.com/Finnhub-Stock-API/finnhub-go"

type ForexRepo struct {
	forex finnhub.Forexrates
}

func NewForexRepository(forex finnhub.Forexrates) *ForexRepo {
	return &ForexRepo{forex: forex}
}

func (fr *ForexRepo) GetBase() string {
	return fr.forex.Base
}

func (fr *ForexRepo) GetQuote() map[string]interface{} {
	return fr.forex.Quote
}
