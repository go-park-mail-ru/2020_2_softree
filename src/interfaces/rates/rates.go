package rates

import (
	"server/src/application"
	"server/src/infrastructure/log"
)

type Rates struct {
	rateApp application.RateApp
	log     log.LogHandler
}

func NewRates(rApp application.RateApp, log log.LogHandler) *Rates {
	return &Rates{rateApp: rApp, log: log}
}
