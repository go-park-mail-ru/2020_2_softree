package authorization

import (
	"github.com/microcosm-cc/bluemonday"
	"server/src/application"
	"server/src/infrastructure/log"
)

type Authentication struct {
	userApp   application.UserApp
	auth      application.UserAuth
	log       log.LogHandler
	sanitizer bluemonday.Policy
}

func NewAuthenticate(uApp application.UserApp, auth application.UserAuth, log log.LogHandler) *Authentication {
	return &Authentication{userApp: uApp, auth: auth, log: log, sanitizer: *bluemonday.UGCPolicy()}
}
