package authorization

import (
	"server/src/application"
	"server/src/infrastructure/log"
)

type Authentication struct {
	userApp application.UserApp
	auth    application.UserAuth
	log     log.LogHandler
}

func NewAuthenticate(uApp application.UserApp, auth application.UserAuth, log log.LogHandler) *Authentication {
	return &Authentication{userApp: uApp, auth: auth, log: log}
}
