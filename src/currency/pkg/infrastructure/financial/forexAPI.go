package financial

import (
	"context"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/spf13/viper"
	"server/src/currency/pkg/domain"
)

type ForexAPI struct {
}

func NewForexAPI() *ForexAPI{
	return &ForexAPI{}
}

func (f *ForexAPI) GetCurrencies() domain.FinancialRepository {
	api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: viper.GetString("finnhub-api.token"),
	})

	forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("USD")})
	finance := NewForexRepository(forexRates)

	return finance
}
