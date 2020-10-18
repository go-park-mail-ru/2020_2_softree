package pureArchAuth

import (
	"net/http"
	"server/src/application"
	"server/src/infrastructure/auth"
)

type Authenticate struct {
	userApp application.UserAppInterface
	auth    auth.AuthInterface
}

func NewAuthenticate (uApp application.UserAppInterface, auth auth.AuthInterface) *Authenticate {
	return &Authenticate{userApp: uApp, auth: auth}
}

func (a *Authenticate) Login(w http.ResponseWriter, r *http.Request) {
}

// Signup?

func (a *Authenticate) Logout(w http.ResponseWriter, r *http.Request) {
}
