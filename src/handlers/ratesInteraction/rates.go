package ratesInteraction

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity/rates"
)

func Rates(w http.ResponseWriter, r *http.Request) {
	result, _ := json.Marshal(rates.Currencies)

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
