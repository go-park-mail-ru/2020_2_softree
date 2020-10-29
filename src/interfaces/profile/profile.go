package profile

import (
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
)

type Profile struct {
	userApp application.UserApp
	auth    auth.AuthHandler
	cookie  auth.TokenHandler
	log     log.LogHandler
}

func NewProfile(
	uApp application.UserApp, auth auth.AuthHandler, cookie auth.TokenHandler, log log.LogHandler) *Profile {
	return &Profile{userApp: uApp, auth: auth, cookie: cookie, log: log}
}
