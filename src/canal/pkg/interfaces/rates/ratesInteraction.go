package rates

import (
	"context"
	"encoding/json"
	"net/http"
	"server/src/canal/pkg/infrastructure/financial"
	"server/src/canal/pkg/infrastructure/persistence"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (rates *Rates) GetRatesFromApi() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: viper.GetString("finnhub-api.token"),
	})

	forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("USD")})
	finance := financial.NewForexRepository(forexRates)

	for range ticker.C {
		/*// exchange works 10:00-20:00
		if time.Now().Hour() > 20 || time.Now().Hour() < 10 {
			continue
		}*/
		forexRates, _, _ = api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("USD")})
		finance = financial.NewForexRepository(forexRates)

		err := rates.rateApp.SaveRates(finance)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"function": "GetRatesFromApi",
			}).Error(err)
			return
		}

		if currencies, err := rates.rateApp.GetInitialCurrency(); len(currencies) == 0 {
			if err = rates.rateApp.SaveCurrency(finance); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "GetRatesFromApi",
				}).Error(err)
				return
			}
		}
		if time.Now().Hour() == 10 && time.Now().Minute() == 2 { // 10:02
			if err = rates.rateApp.SaveCurrency(finance); err != nil {
				logrus.WithFields(logrus.Fields{
					"function": "GetRatesFromApi",
				}).Error(err)
				return
			}
		}
	}
}

func (rates *Rates) GetRates(w http.ResponseWriter, r *http.Request) {
	resRates, err := rates.rateApp.GetRates()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetRates",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRates",
		}).Error(err)
	}
}

func (rates *Rates) GetURLRate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	title := vars["title"]
	if !validate(title) {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusBadRequest,
			"function": "GetURLRate",
			"title":    title,
		}).Error("Bad title")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resRates, err := rates.rateApp.GetRate(title)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetURLRate",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetURLRate",
		}).Error(err)
	}
}

func validate(title string) bool {
	lenOfCurrency := 3
	if len(title) != lenOfCurrency {
		return false
	}

	for _, rate := range persistence.ListOfCurrencies {
		if rate == title {
			return true
		}
	}

	return false
}

func (rates *Rates) GetMarkets(w http.ResponseWriter, r *http.Request) {
	type curr struct {
		Base  string `json:"base"`
		Title string `json:"title"`
	}
	resRates := [...]curr{
		{Base: "USD", Title: "EUR"},
		{Base: "USD", Title: "RUB"},
		{Base: "USD", Title: "JPY"},
		{Base: "USD", Title: "GBP"},
		{Base: "RUB", Title: "ZAR"},
		{Base: "RUB", Title: "EUR"},
		{Base: "RUB", Title: "BRL"},
		{Base: "RUB", Title: "ILS"},
	}
	result, _ := json.Marshal(resRates)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetMarkets",
		}).Error(err)
	}
}
