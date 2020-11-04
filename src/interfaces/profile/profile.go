package profile

import (
	"server/src/application"
	"server/src/infrastructure/log"
)

type Profile struct {
	userApp application.UserApp
	auth    application.UserAuth
	log     log.LogHandler
}

func NewProfile(uApp application.UserApp, auth application.UserAuth, log log.LogHandler) *Profile {
	return &Profile{userApp: uApp, auth: auth, log: log}
}
