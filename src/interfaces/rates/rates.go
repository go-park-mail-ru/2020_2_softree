package rates

import (
	"server/src/application"
)

type Rates struct {
	rateApp application.RateApp
}

func NewRates(rApp application.RateApp) *Rates {
	return &Rates{rateApp: rApp}
}
