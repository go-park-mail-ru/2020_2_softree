package profile

import (
	"server/src/application"
	"server/src/infrastructure/auth"
)

type Profile struct {
	userApp application.UserAppInterface
	auth    auth.AuthInterface
	cookie  auth.CookieInterface
}

func NewProfile(
	uApp application.UserAppInterface, auth auth.AuthInterface, cookie auth.CookieInterface) *Profile {
	return &Profile{userApp: uApp, auth: auth, cookie: cookie}
}
