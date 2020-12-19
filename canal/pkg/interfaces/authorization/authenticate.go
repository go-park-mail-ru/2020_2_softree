package authorization

import (
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/logger"
)

type Authentication struct {
	profileLogic repository.ProfileLogic
	authLogic    repository.AuthLogic
	logger       logger.Logrus
}

func NewAuthentication(profileLogic repository.ProfileLogic, authLogic repository.AuthLogic) *Authentication {
	return &Authentication{profileLogic, authLogic, *logger.NewLogrus()}
}
