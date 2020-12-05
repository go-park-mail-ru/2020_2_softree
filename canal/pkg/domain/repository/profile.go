package repository

import (
	"context"
	"server/canal/pkg/domain/entity"
)

type ProfileLogic interface {
	UpdateAvatar(ctx context.Context, user entity.User) (entity.Description, entity.PublicUser)
}
