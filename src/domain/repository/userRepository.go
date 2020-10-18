package repository

import (
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, *jsonRealisation.ErrorJSON)
	UpdateUser(uint64, entity.User) (entity.User, error)
	DeleteUser(uint64) error
	GetUser(uint64) (*entity.User, error)
}
