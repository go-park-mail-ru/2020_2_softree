package profile

import (
	"server/src/application"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
)

type Profile struct {
	userApp application.UserAppHandler
	auth    auth.AuthHandler
	cookie  auth.CookieHandler
	log     log.LogHandler
}

func NewProfile(
	uApp application.UserAppHandler, auth auth.AuthHandler, cookie auth.CookieHandler, log log.LogHandler) *Profile {
	return &Profile{userApp: uApp, auth: auth, cookie: cookie, log: log}
}
