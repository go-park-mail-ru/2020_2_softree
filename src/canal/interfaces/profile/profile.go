package profile

import (
	"server/src/canal/application"
	profile "server/src/profile/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	userApp   profile.ProfileServiceClient
	rateApp   application.RateApp
	auth      application.UserAuth
	sanitizer bluemonday.Policy
}

func NewProfile(
	uApp profile.ProfileServiceClient, auth application.UserAuth, rApp application.RateApp) *Profile {
	return &Profile{userApp: uApp, auth: auth, sanitizer: *bluemonday.UGCPolicy(), rateApp: rApp}
}
