package authorization

import (
	session "server/src/authorization/pkg/session/gen"
	"server/src/canal/pkg/domain/repository"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Authentication struct {
	profile   profile.ProfileServiceClient
	auth      session.AuthorizationServiceClient
	security  repository.Utils
	sanitizer bluemonday.Policy
}

func NewAuthenticate(profile profile.ProfileServiceClient,
	auth session.AuthorizationServiceClient, security repository.Utils) *Authentication {
	return &Authentication{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy(), security: security}
}
