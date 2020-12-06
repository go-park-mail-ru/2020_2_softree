package profile

import (
	"github.com/prometheus/client_golang/prometheus"
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/logger"
	"server/canal/pkg/infrastructure/metric"
	currency "server/currency/pkg/currency/gen"
)

type Profile struct {
	profileLogic repository.ProfileLogic
	paymentLogic repository.PaymentLogic
	logger       logger.Logrus
	rates        currency.CurrencyServiceClient
	Hits         prometheus.CounterVec
}

func NewProfile(currency currency.CurrencyServiceClient) *Profile {
	return &Profile{
		rates:     currency,
		Hits:      *metric.Metric,
	}
}
