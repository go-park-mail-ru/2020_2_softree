package repository

import (
	"server/src/canal/pkg/domain/entity"
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
	UpdateUserAvatar(int64, entity.User) error
}

type userPasswordUpdater interface {
	UpdateUserPassword(int64, entity.User) error
}

type userEraser interface {
	DeleteUser(int64) error
}

type userReceiverById interface {
	GetUserById(int64) (entity.User, error)
}

type userReceiverByFormData interface {
	GetUserByLogin(string, string) (entity.User, error)
}

type userReceiverWatchlist interface {
	GetUserWatchlist(int64) ([]entity.Currency, error)
}

type userCheckExistence interface {
	CheckExistence(string) (bool, error)
}

type userCheckPassword interface {
	CheckPassword(int64, string) (bool, error)
}
