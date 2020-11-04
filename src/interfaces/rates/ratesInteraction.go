package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"net/http"
	"server/src/domain/entity"
)

func (rates *Rates) GetRates(w http.ResponseWriter, r *http.Request) {
	api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "bttn28748v6ojt2hev60",
	})

	var rate entity.Rate
	forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("RUB")})
	for name, quote := range forexRates.Quote {
		rate.Base = forexRates.Base
		rate.Currency = name
		rate.Value = fmt.Sprintf("%.6f", quote.(float64))

		var err error
		if rate, err = rates.rateApp.SaveRate(rate); err != nil {
			rates.log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	resRates, err := rates.rateApp.GetRates()
	if err != nil {
		rates.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}
