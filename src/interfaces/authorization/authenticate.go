package authorization

import (
	"server/src/application"

	"github.com/microcosm-cc/bluemonday"
)

type Authentication struct {
	userApp   application.UserApp
	auth      application.UserAuth
	sanitizer bluemonday.Policy
}

func NewAuthenticate(uApp application.UserApp, auth application.UserAuth) *Authentication {
	return &Authentication{userApp: uApp, auth: auth, sanitizer: *bluemonday.UGCPolicy()}
}
