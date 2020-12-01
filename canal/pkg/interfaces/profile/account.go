package profile

import (
	"github.com/prometheus/client_golang/prometheus"
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
	Hits      prometheus.CounterVec
}

func NewProfile(
	profile profile.ProfileServiceClient, security repository.Utils, currency currency.CurrencyServiceClient) *Profile {
	return &Profile{
		profile:   profile,
		rates:     currency,
		security:  security,
		sanitizer: *bluemonday.UGCPolicy(),
		Hits:      *prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hits"}, []string{"status"}),
	}
}
