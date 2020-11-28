package profile

import (
	"server/canal/pkg/domain/repository"
	currency "server/currency/pkg/currency/gen"
	profile "server/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Profile struct {
	profile   profile.ProfileServiceClient
	rates     currency.CurrencyServiceClient
	security  repository.Utils
	sanitizer bluemonday.Policy
}

func NewProfile(
	profile profile.ProfileServiceClient, security repository.Utils, currency currency.CurrencyServiceClient) *Profile {
	return &Profile{profile: profile, sanitizer: *bluemonday.UGCPolicy(), rates: currency, security: security}
}
