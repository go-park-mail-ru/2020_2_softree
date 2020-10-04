package ratesInteraction

import (
	"encoding/json"
	"net/http"
	"server/domain/entity"
	"strconv"
)

func Rates(w http.ResponseWriter, r *http.Request) {
	result := make([]byte, 0)
	id, err := strconv.Atoi(r.FormValue("id"))
	if id == 0 {
		var quotations entity.Quotations
		quotations.Add(
			entity.CurrencyQuotation{Value: 1234, Title: "title1", Change: 12},
			entity.CurrencyQuotation{Value: 12, Title: "title2", Change: 1},
		)

		result, err = json.Marshal(quotations)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		result, err = json.Marshal(entity.FindById(uint64(id)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
