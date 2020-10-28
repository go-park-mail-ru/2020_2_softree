package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type userApp struct {
	ur repository.UserRepository
}

type UserAppHandler interface {
	SaveUser(entity.User) (entity.User, error)
	UpdateUser(uint64, entity.User) (entity.User, error)
	DeleteUser(uint64) error
	GetUser(uint64) (*entity.User, error)
}

func (ua *userApp) SaveUser(u entity.User) (entity.User, error) {
	return ua.ur.SaveUser(u)
}

func (ua *userApp) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	return ua.ur.UpdateUser(id, u)
}

func (ua *userApp) DeleteUser(id uint64) error {
	return ua.ur.DeleteUser(id)
}

func (ua *userApp) GetUser(id uint64) (*entity.User, error) {
	return ua.ur.GetUser(id)
}
