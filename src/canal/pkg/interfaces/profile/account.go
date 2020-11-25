package profile

import (
	"server/src/canal/pkg/application"
	"server/src/canal/pkg/domain/repository"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	profile   profile.ProfileServiceClient
	rateApp   application.RateApp
	security  repository.Utils
	sanitizer bluemonday.Policy
}

func NewProfile(
	profile profile.ProfileServiceClient, rApp application.RateApp, security repository.Utils) *Profile {
	return &Profile{profile: profile, sanitizer: *bluemonday.UGCPolicy(), rateApp: rApp, security: security}
}
