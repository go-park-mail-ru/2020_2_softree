package pureArchAuth

import (
	"net/http"
	"server/src/application"
	"server/src/infrastructure/auth"
)

type Authenticate struct {
	userApp application.UserAppInterface
	auth    auth.AuthInterface
	cookie  auth.CookieInterface
}

func NewAuthenticate(
	uApp application.UserAppInterface, auth auth.AuthInterface, cookie auth.CookieInterface) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth, cookie: cookie}
}

func (a *Authenticate) Login(w http.ResponseWriter, r *http.Request) {
}

func (a *Authenticate) Logout(w http.ResponseWriter, r *http.Request) {
}
