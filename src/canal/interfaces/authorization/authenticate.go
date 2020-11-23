package authorization

import (
	session "server/src/authorization/session/gen"
	profile "server/src/profile/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Authentication struct {
	profile   profile.ProfileServiceClient
	auth      session.AuthorizationServiceClient
	sanitizer bluemonday.Policy
}

func NewAuthenticate(profile profile.ProfileServiceClient, auth session.AuthorizationServiceClient) *Authentication {
	return &Authentication{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy()}
}
