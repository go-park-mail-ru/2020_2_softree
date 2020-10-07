package ratesInteraction

import (
	"encoding/json"
	"net/http"
	"server/domain/entity/rates"
)

func Rates(w http.ResponseWriter, r *http.Request) {
	result, err := json.Marshal(rates.Currencies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
