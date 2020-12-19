package rates

import (
	json "github.com/mailru/easyjson"
	"net/http"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/metric"
	"time"
)

func (rates *Rates) GetAllLatestRates(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetAllLatestRates")

	desc, out, err := rates.currencyApp.GetAllLatestCurrencies(r)
	if err != nil {
		rates.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	result, err := json.Marshal(out)
	if err != nil {
		desc = entity.Description{Function: "GetAllLatestRates", Action: "Marshal", Status: http.StatusInternalServerError}
		rates.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err = w.Write(result); err != nil {
		rates.logger.Error(entity.Description{Function: "GetAllLatestRates", Action: "Write"}, err)
	}
}

func (rates *Rates) GetURLRate(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetURLRate")

	desc, out, err := rates.currencyApp.GetURLCurrencies(r)
	if err != nil {
		rates.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	result, err := json.Marshal(out)
	if err != nil {
		desc = entity.Description{Function: "GetURLRate", Action: "Marshal", Status: http.StatusInternalServerError}
		rates.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err = w.Write(result); err != nil {
		rates.logger.Error(entity.Description{Function: "GetURLRate", Action: "Write"}, err)
	}
}

func (rates *Rates) GetMarkets(w http.ResponseWriter, r *http.Request) {
	defer metric.RecordTimeMetric(time.Now(), "GetMarkets")

	desc, out, err := rates.currencyApp.GetMarkets()
	if err != nil {
		rates.logger.Error(desc, err)
		w.WriteHeader(desc.Status)

		metric.RecordHitMetric(desc.Status, r.URL.Path)
		return
	}

	result, err := json.Marshal(out)
	if err != nil {
		desc = entity.Description{Function: "GetMarkets", Action: "Marshal", Status: http.StatusInternalServerError}
		rates.logger.Error(desc, err)
		w.WriteHeader(http.StatusInternalServerError)

		metric.RecordHitMetric(http.StatusInternalServerError, r.URL.Path)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	metric.RecordHitMetric(http.StatusOK, r.URL.Path)
	if _, err = w.Write(result); err != nil {
		rates.logger.Error(entity.Description{Function: "GetMarkets", Action: "Write"}, err)
	}
}
