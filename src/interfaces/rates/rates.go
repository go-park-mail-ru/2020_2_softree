package rates

import (
	"server/src/application"
	"server/src/infrastructure/log"
)

type Rates struct {
	rateApp application.RateApp
	auth    application.UserAuth
	log     log.LogHandler
}

func NewRates(rApp application.RateApp, auth application.UserAuth, log log.LogHandler) *Rates {
	return &Rates{rateApp: rApp, auth: auth, log: log}
}
