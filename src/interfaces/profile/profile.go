package profile

import (
	"github.com/microcosm-cc/bluemonday"
	"server/src/application"
	"server/src/infrastructure/log"
)

type Profile struct {
	userApp   application.UserApp
	auth      application.UserAuth
	log       log.LogHandler
	sanitizer bluemonday.Policy
}

func NewProfile(uApp application.UserApp, auth application.UserAuth, log log.LogHandler) *Profile {
	return &Profile{userApp: uApp, auth: auth, log: log, sanitizer: *bluemonday.UGCPolicy()}
}
