package pureArchAuth

import (
	"server/src/application"
	"server/src/infrastructure/auth"
)

type Authenticate struct {
	userApp application.UserAppInterface
	auth    auth.AuthInterface
	cookie  auth.CookieInterface
	// logger will added later
}

func NewAuthenticate(
	uApp application.UserAppInterface, auth auth.AuthInterface, cookie auth.CookieInterface) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth, cookie: cookie}
}
