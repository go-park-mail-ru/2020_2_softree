package repository

import (
	"context"
	"net/http"
	"server/canal/pkg/domain/entity"
)

type AuthLogic interface {
	Login(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, http.Cookie, error)
	Logout(ctx context.Context, cookie *http.Cookie) (entity.Description, http.Cookie, error)
	Signup(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, http.Cookie, error)
	Authenticate(ctx context.Context, userId int64) (entity.Description, entity.PublicUser, error)
	Auth(ctx context.Context, cookie *http.Cookie) (entity.Description, int64, error)
}
