package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type UserApp struct {
	ur repository.UserRepository
}

func NewUserApp(repo repository.UserRepository) *UserApp {
	return &UserApp{ur: repo}
}

func (ua *UserApp) SaveUser(u entity.User) (entity.User, error) {
	return ua.ur.SaveUser(u)
}

func (ua *UserApp) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	return ua.ur.UpdateUser(id, u)
}

func (ua *UserApp) DeleteUser(id uint64) error {
	return ua.ur.DeleteUser(id)
}

func (ua *UserApp) GetUser(id uint64) (entity.User, error) {
	return ua.ur.GetUser(id)
}
