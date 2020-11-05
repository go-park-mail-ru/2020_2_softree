package rates

import (
	"context"
	"encoding/json"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"net/http"
	"server/src/domain/repository"
	"server/src/infrastructure/financial"
)

func (rates *Rates) ForexRates(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
		auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
			Key: "bttn28748v6ojt2hev60",
		})

		forexRates, _, _ := api.ForexRates(auth, &finnhub.ForexRatesOpts{Base: optional.NewString("USD")})
		finance := financial.NewForexRepository(forexRates)

		ctx := context.WithValue(r.Context(), "finance", finance)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	}
}

func (rates *Rates) GetRates(w http.ResponseWriter, r *http.Request) {
	finance := r.Context().Value("finance").(repository.FinancialRepository)

	resRates, err := rates.rateApp.SaveRates(finance)
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
