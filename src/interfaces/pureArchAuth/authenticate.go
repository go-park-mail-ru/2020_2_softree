package pureArchAuth

import (
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
)

type Authenticate struct {
	userApp application.UserApp
	auth    auth.AuthHandler
	cookie  auth.TokenHandler
	log     log.LogHandler
}

func NewAuthenticate(
	uApp application.UserApp, auth auth.AuthHandler, cookie auth.TokenHandler, log log.LogHandler) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth, cookie: cookie, log: log}
}
