package rates

import (
	"server/canal/pkg/application"
	"server/canal/pkg/infrastructure/logger"
)

type Rates struct {
	currencyApp application.CurrencyApp
	logger      logger.Logrus
}

func NewRates(currencyApp application.CurrencyApp) *Rates {
	return &Rates{currencyApp, *logger.NewLogrus()}
}
