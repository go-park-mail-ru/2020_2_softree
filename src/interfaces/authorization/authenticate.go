package authorization

import (
	"server/src/application"
	session "server/src/authorizationService/session/gen"

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
