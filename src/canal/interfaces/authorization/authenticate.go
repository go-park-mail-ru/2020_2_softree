package authorization

import (
	session "server/src/authorization/session/gen"
	"server/src/canal/application"

	"github.com/microcosm-cc/bluemonday"
)

type Authentication struct {
	userApp   application.UserApp
	auth      session.AuthorizationServiceClient
	sanitizer bluemonday.Policy
}

func NewAuthenticate(uApp application.UserApp, auth session.AuthorizationServiceClient) *Authentication {
	return &Authentication{userApp: uApp, auth: auth, sanitizer: *bluemonday.UGCPolicy()}
}
