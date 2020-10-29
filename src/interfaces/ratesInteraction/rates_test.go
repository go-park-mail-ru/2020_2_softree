package ratesInteraction

import (
	"encoding/json"
	"net/http/httptest"
	"server/src/domain/entity/rates"
	"testing"
)

func TestRates(t *testing.T) {
	url := "http://127.0.0.1:8000/rates"

	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	Rates(w, req)

	var currency []rates.Currency
	decoder := json.NewDecoder(w.Body)
	err := decoder.Decode(&currency)
	if err != nil {
		t.Fatal(err)
	}

	if len(currency) == 0 {
		t.Fatal("fail to get currencies")
	}
}
