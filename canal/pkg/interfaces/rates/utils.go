package rates

import (
	"server/currency/pkg/infrastructure/persistence"
	"strconv"
)

func (rates *Rates) recordHitMetric(code int) {
	rates.Hits.WithLabelValues(strconv.Itoa(code)).Inc()
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
