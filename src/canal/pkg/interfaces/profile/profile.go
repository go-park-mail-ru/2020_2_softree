package profile

import (
	session "server/src/authorization/pkg/session/gen"
	"server/src/canal/pkg/application"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	profile   profile.ProfileServiceClient
	rateApp   application.RateApp
	auth      session.AuthorizationServiceClient
	sanitizer bluemonday.Policy
}

func NewProfile(
	profile profile.ProfileServiceClient, auth session.AuthorizationServiceClient, rApp application.RateApp) *Profile {
	return &Profile{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy(), rateApp: rApp}
}
