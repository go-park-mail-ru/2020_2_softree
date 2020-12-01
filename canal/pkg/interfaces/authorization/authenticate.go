package authorization

import (
	"github.com/prometheus/client_golang/prometheus"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/repository"
	profile "server/profile/pkg/profile/gen"

	"github.com/microcosm-cc/bluemonday"
)

type Authentication struct {
	profile   profile.ProfileServiceClient
	auth      session.AuthorizationServiceClient
	security  repository.Utils
	sanitizer bluemonday.Policy
	Hits      prometheus.CounterVec
}

func NewAuthenticate(profile profile.ProfileServiceClient,
	auth session.AuthorizationServiceClient, security repository.Utils) *Authentication {
	return &Authentication{
		profile:   profile,
		auth:      auth,
		security:  security,
		sanitizer: *bluemonday.UGCPolicy(),
		Hits:      *prometheus.NewCounterVec(prometheus.CounterOpts{Name: "hits"}, []string{"status"}),
	}
}
