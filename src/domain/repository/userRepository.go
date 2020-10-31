package repository

import (
	"server/src/domain/entity"
)

type UserRepository interface {
	saver
	updater
	eraser
	receiverById
	receiverByFormData
}

type saver interface {
	SaveUser(entity.User) (entity.User, error)
}

type updater interface {
	UpdateUser(uint64, entity.User) (entity.User, error)
}

type eraser interface {
	DeleteUser(uint64) error
}

type receiverById interface {
	GetUserById(uint64) (entity.User, error)
}

type receiverByFormData interface {
	GetUserByLogin(string, string) (entity.User, error)
}
