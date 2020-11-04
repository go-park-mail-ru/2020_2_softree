package authorization

import (
	"server/src/application"
	"server/src/infrastructure/log"
)

type Authenticate struct {
	userApp application.UserApp
	auth    application.UserAuth
	log     log.LogHandler
}

func NewAuthenticate(uApp application.UserApp, auth application.UserAuth, log log.LogHandler) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth, log: log}
}
