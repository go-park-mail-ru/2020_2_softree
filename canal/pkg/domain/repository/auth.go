package repository

import (
	"context"
	"net/http"
	authorization "server/authorization/pkg/session/gen"
	"server/canal/pkg/domain/entity"
)

type AuthLogic interface {
	Login(ctx context.Context, user entity.User)  (entity.Description, entity.PublicUser, http.Cookie, error)
	Logout(ctx context.Context, session authorization.Session) (entity.Description, error)
	Authenticate(ctx context.Context, ) (entity.Description, error)
}