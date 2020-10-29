package repository

import (
	"server/src/domain/entity"
)

type UserRepository interface {
	SaveUser(entity.User) (entity.User, error)
	UpdateUser(uint64, entity.User) (entity.User, error)
	DeleteUser(uint64) error
	GetUser(uint64) (entity.User, error)
}
