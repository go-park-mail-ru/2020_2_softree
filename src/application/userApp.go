package application

import (
	"server/src/domain/entity"
	"server/src/domain/repository"
)

type UserApp struct {
	userRepository repository.UserRepository
}

func NewUserApp(repo repository.UserRepository) *UserApp {
	return &UserApp{userRepository: repo}
}

func (ua *UserApp) SaveUser(u entity.User) (entity.User, error) {
	return ua.userRepository.SaveUser(u)
}

func (ua *UserApp) UpdateUserAvatar(id uint64, u entity.User) error {
	return ua.userRepository.UpdateUserAvatar(id, u)
}

func (ua *UserApp) UpdateUserPassword(id uint64, u entity.User) error {
	return ua.userRepository.UpdateUserPassword(id, u)
}

func (ua *UserApp) DeleteUser(id uint64) error {
	return ua.userRepository.DeleteUser(id)
}

func (ua *UserApp) GetUserById(id uint64) (entity.User, error) {
	return ua.userRepository.GetUserById(id)
}

func (ua *UserApp) GetUserByLogin(email, password string) (entity.User, error) {
	return ua.userRepository.GetUserByLogin(email, password)
}

func (ua *UserApp) GetUserWatchlist(id uint64) ([]entity.Currency, error) {
	return ua.userRepository.GetUserWatchlist(id)
}

func (ua *UserApp) CheckExistence(email string) (bool, error) {
	return ua.userRepository.CheckExistence(email)
}

func (ua *UserApp) CheckPassword(id uint64, password string) (bool, error) {
	return ua.userRepository.CheckPassword(id, password)
}
