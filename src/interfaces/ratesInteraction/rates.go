package ratesInteraction

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"net/http"
	"server/src/domain/entity"
)

func Rates(w http.ResponseWriter, r *http.Request) {
	api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "bttn28748v6ojt2hev60",
	})

	var rate entity.Rate
	var rates entity.Rates
	forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("RUB")})
	for name, quote := range forexRates.Quote {
		rate.Name = fmt.Sprintf("%s/%s", forexRates.Base, name)
		rate.Value = fmt.Sprintf("%.6f", quote.(float64))

		rates.Values = append(rates.Values, rate)
	}

	result, _ := json.Marshal(rates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}
