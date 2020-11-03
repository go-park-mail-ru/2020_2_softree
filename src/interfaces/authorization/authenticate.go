package authorization

import (
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
)

type Authenticate struct {
	userApp application.UserApp
	auth    application.UserAuth
	cookie  auth.TokenHandler
	log     log.LogHandler
}

func NewAuthenticate(
	uApp application.UserApp, auth application.UserAuth, cookie auth.TokenHandler, log log.LogHandler) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth, cookie: cookie, log: log}
}
