package rates

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	currency "server/currency/pkg/currency/gen"
)

func (rates *Rates) GetRates(w http.ResponseWriter, r *http.Request) {
	resRates, err := rates.currencyService.GetRates(r.Context(), &currency.Empty{})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetRates",
			"action":   "GetRates",
		}).Error(err)

		rates.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(resRates.Rates)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetRates",
			"action":   "Marshal",
		}).Error(err)

		rates.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rates.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetRates",
			"action":   "Write",
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
			"action":   "validate",
			"title":    title,
		}).Error("Bad title")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resRates, err := rates.currencyService.GetRate(r.Context(), &currency.CurrencyTitle{Title: title})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetURLRate",
			"action":   "GetRate",
			"title":    title,
		}).Error(err)

		rates.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	result, err := json.Marshal(resRates.Rates)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetURLRate",
			"action":   "Marshal",
		}).Error(err)

		rates.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rates.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetURLRate",
			"action":   "Write",
		}).Error(err)
	}
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
		{Base: "RUB", Title: "GBP"},
		{Base: "RUB", Title: "EUR"},
		{Base: "RUB", Title: "BRL"},
		{Base: "RUB", Title: "ILS"},
		{Base: "RUB", Title: "JPY"},
	}
	result, err := json.Marshal(resRates)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"status":   http.StatusInternalServerError,
			"function": "GetMarkets",
			"action":   "Marshal",
		}).Error(err)

		rates.recordHitMetric(http.StatusInternalServerError)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rates.recordHitMetric(http.StatusOK)
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(result); err != nil {
		logrus.WithFields(logrus.Fields{
			"function": "GetMarkets",
			"action":   "Write",
		}).Error(err)
	}
}
