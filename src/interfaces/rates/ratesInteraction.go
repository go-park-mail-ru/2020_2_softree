package rates

import (
	"context"
	"encoding/json"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/gorilla/mux"
	"net/http"
	"server/src/domain/entity"
	"server/src/infrastructure/financial"
	"time"
)

func (rates *Rates) GetRatesFromApi() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: "bttn28748v6ojt2hev60",
	})

	for range ticker.C {
		forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("USD")})
		finance := financial.NewForexRepository(forexRates)

		err := rates.rateApp.SaveRates(finance)
		if err != nil {
			rates.log.Print(err)
			return
		}

		if time.Now().Hour() == 10 && time.Now().Minute() == 2 {  // 10:02
			if err = rates.rateApp.SaveCurrency(finance);  err != nil {
				rates.log.Print(err)
				return
			}
		}
	}
}

func (rates *Rates) GetRates(w http.ResponseWriter, r *http.Request) {
	resRates, err := rates.rateApp.GetRates()
	if err != nil {
		rates.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}

func (rates *Rates) GetURLRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	resRates, err := rates.rateApp.GetRate(vars["title"])
	if err != nil {
		rates.log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}

func (rates *Rates) GetMarkets(w http.ResponseWriter, r *http.Request) {
	resRates := entity.Currency{Base: "USD", Title: "EUR"}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(result)
}
