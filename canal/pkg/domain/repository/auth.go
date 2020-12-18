package repository

import (
	"context"
	authorization "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/entity"
	profile "server/profile/pkg/profile/gen"
)

type AuthLogic interface {
	Login(ctx context.Context, user profile.User) (entity.Description, error)
	Logout(ctx context.Context, session authorization.Session) (entity.Description, error)
	Authenticate(ctx context.Context, ) (entity.Description, error)
}