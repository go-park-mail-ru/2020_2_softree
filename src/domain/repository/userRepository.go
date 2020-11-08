package repository

import (
	"server/src/domain/entity"
)

type UserRepository interface {
	userSaver
	userUpdater
	userEraser
	userReceiverById
	userReceiverByFormData
	userReceiverWatchlist
}

type userSaver interface {
	SaveUser(entity.User) (entity.User, error)
}

type userUpdater interface {
	UpdateUser(uint64, entity.User) (entity.User, error)
}

type userEraser interface {
	DeleteUser(uint64) error
}

type userReceiverById interface {
	GetUserById(uint64) (entity.User, error)
}

type userReceiverByFormData interface {
	GetUserByLogin(string, string) (entity.User, error)
}

type userReceiverWatchlist interface {
	GetUserWatchlist(uint64) ([]entity.Currency, error)
}
