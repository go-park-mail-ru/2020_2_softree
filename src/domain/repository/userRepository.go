package repository

import (
	"server/src/domain/entity"
)

type UserRepository interface {
	userSaver
	userAvatarUpdater
	userPasswordUpdater
	userEraser
	userReceiverById
	userReceiverByFormData
	userReceiverWatchlist
	userCheckExistence
	userCheckPassword
}

type userSaver interface {
	SaveUser(entity.User) (entity.User, error)
}

type userAvatarUpdater interface {
	UpdateUserAvatar(uint64, entity.User) error
}

type userPasswordUpdater interface {
	UpdateUserPassword(uint64, entity.User) error
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

type userCheckExistence interface {
	CheckExistence(string) (bool, error)
}

type userCheckPassword interface {
	CheckPassword(uint64, string) (bool, error)
}
