package profile

import (
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/logger"
)

type Profile struct {
	profileLogic repository.ProfileLogic
	paymentLogic repository.PaymentLogic
	logger       logger.Logrus
}

func NewProfile(profileLogic repository.ProfileLogic, paymentLogic repository.PaymentLogic) *Profile {
	return &Profile{profileLogic, paymentLogic, *logger.NewLogrus()}
}
