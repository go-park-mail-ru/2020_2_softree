package application

import (
	"net/http"
	"server/src/domain/repository"
)

type UserAuth struct {
	authRepository repository.AuthRepository
}

func NewUserAuth(auth repository.AuthRepository) *UserAuth {
	return &UserAuth{authRepository: auth}
}

func (ua *UserAuth) CreateAuth(id uint64) (http.Cookie, error) {
	return ua.authRepository.CreateAuth(id)
}

func (ua *UserAuth) CheckAuth(sessionValue string) (uint64, error) {
	return ua.authRepository.CheckAuth(sessionValue)
}

func (ua *UserAuth) DeleteAuth(value string) error {
	return ua.authRepository.DeleteAuth(value)
}

func (ua *UserAuth) CreateCookie() (http.Cookie, error) {
	return ua.authRepository.CreateCookie()
}
