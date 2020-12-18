package application

import (
	"context"
	"github.com/microcosm-cc/bluemonday"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/domain/repository"
	profile "server/profile/pkg/profile/gen"
	authorization "server/authorization/pkg/session/gen"
)

type AuthApp struct {
	profile profile.ProfileServiceClient
	auth authorization.AuthorizationServiceClient
	sanitizer bluemonday.Policy
	security  repository.Utils
}

func NewAuthApp(profile profile.ProfileServiceClient, auth authorization.AuthorizationServiceClient, security repository.Utils) *AuthApp {
	return &AuthApp{profile: profile, auth: auth, sanitizer: *bluemonday.UGCPolicy(), security: security}
}

func (authApp *AuthApp) Login(ctx context.Context, user profile.User) (entity.Description, error) {

}

func (authApp *AuthApp) Logout(ctx context.Context, session authorization.Session) (entity.Description, error) {

}

func (authApp *AuthApp) Authenticate(ctx context.Context, userId authorization.UserID) (entity.Description, error) {

}
