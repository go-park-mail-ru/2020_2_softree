package profile

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/sirupsen/logrus"
	"server/src/application"
)

type Profile struct {
	userApp   application.UserApp
	auth      application.UserAuth
	log       logrus.Logger
	sanitizer bluemonday.Policy
}

func NewProfile(uApp application.UserApp, auth application.UserAuth, log *logrus.Logger) *Profile {
	return &Profile{userApp: uApp, auth: auth, log: *log, sanitizer: *bluemonday.UGCPolicy()}
}
