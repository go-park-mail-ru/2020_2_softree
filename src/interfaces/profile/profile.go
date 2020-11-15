package profile

import (
	"server/src/application"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	userApp   application.UserApp
	rateApp   application.RateApp
	auth      application.UserAuth
	sanitizer bluemonday.Policy
}

func NewProfile(
	uApp application.UserApp, auth application.UserAuth, rApp application.RateApp) *Profile {
	return &Profile{userApp: uApp, auth: auth, sanitizer: *bluemonday.UGCPolicy(), rateApp: rApp}
}
