package repository

import (
	"context"
	"server/canal/pkg/domain/entity"
)

type ProfileLogic interface {
	UpdateAvatar(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, error)
	UpdatePassword(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser, error)
	ReceiveUser(ctx context.Context, id int64) (entity.Description, entity.PublicUser, error)
	ReceiveWatchlist(ctx context.Context, id int64) (entity.Description, entity.Currencies, error)
}
