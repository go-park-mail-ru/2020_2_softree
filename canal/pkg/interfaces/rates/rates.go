package rates

import (
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/logger"
)

type Rates struct {
	currencyApp repository.CurrencyLogic
	logger      logger.Logrus
}

func NewRates(currencyApp repository.CurrencyLogic) *Rates {
	return &Rates{currencyApp, *logger.NewLogrus()}
}
