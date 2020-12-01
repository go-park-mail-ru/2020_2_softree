package authorization

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/prometheus/client_golang/prometheus"
	session "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/repository"
	"server/canal/pkg/infrastructure/metric"
	profile "server/profile/pkg/profile/gen"
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
		Hits:      *metric.Metric,
	}
}
