package profile

import (
	"github.com/prometheus/client_golang/prometheus"
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/logger"
	"server/canal/pkg/infrastructure/metric"
)

type Profile struct {
	profileLogic repository.ProfileLogic
	paymentLogic repository.PaymentLogic
	logger       logger.Logrus
	Hits         prometheus.CounterVec
}

func NewProfile(profileLogic repository.ProfileLogic, paymentLogic repository.PaymentLogic) *Profile {
	return &Profile{profileLogic, paymentLogic, *logger.NewLogrus(), *metric.Metric}
}
