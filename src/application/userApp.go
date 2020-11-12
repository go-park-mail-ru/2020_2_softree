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

func (ua *UserApp) UpdateUserAvatar(id uint64, u entity.User) error {
	return ua.ur.UpdateUserAvatar(id, u)
}

func (ua *UserApp) UpdateUserPassword(id uint64, u entity.User) error {
	return ua.ur.UpdateUserPassword(id, u)
}

func (ua *UserApp) DeleteUser(id uint64) error {
	return ua.ur.DeleteUser(id)
}

func (ua *UserApp) GetUserById(id uint64) (entity.User, error) {
	return ua.ur.GetUserById(id)
}

func (ua *UserApp) GetUserByLogin(email, password string) (entity.User, error) {
	return ua.ur.GetUserByLogin(email, password)
}

func (ua *UserApp) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	return ua.ur.GetUserWatchlist(id)
}

func (ua *UserApp) CheckExistence(email string) (bool, error) {
	return ua.ur.CheckExistence(email)
}

func (ua *UserApp) CheckPassword(id uint64, password string) (bool, error) {
	return ua.ur.CheckPassword(id, password)
}
