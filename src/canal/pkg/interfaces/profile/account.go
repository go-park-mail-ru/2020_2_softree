package profile

import (
	"server/src/canal/pkg/domain/repository"
	currency "server/src/currency/pkg/currency/gen"
	profile "server/src/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	profile   profile.ProfileServiceClient
	currency  currency.CurrencyServiceClient
	security  repository.Utils
	sanitizer bluemonday.Policy
}

func NewProfile(
	profile profile.ProfileServiceClient, security repository.Utils, currency currency.CurrencyServiceClient) *Profile {
	return &Profile{profile: profile, sanitizer: *bluemonday.UGCPolicy(), currency: currency, security: security}
}
