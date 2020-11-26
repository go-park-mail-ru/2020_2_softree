package profile

import (
	session "server/src/authorization/pkg/session/gen"
	currency "server/src/currency/pkg/currency/gen"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	profile   profile.ProfileServiceClient
	currency   currency.CurrencyServiceClient
	auth      session.AuthorizationServiceClient
	sanitizer bluemonday.Policy
}

func NewProfile(
	profile profile.ProfileServiceClient, auth session.AuthorizationServiceClient, currency   currency.CurrencyServiceClient) *Profile {
	return &Profile{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy(), currency: currency}
}
