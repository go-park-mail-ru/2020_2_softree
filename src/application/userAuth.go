package application

import (
	"net/http"
	"server/src/domain/repository"
)

type UserAuth struct {
	ur repository.AuthRepository
}

func NewUserAuth(auth repository.AuthRepository) *UserAuth {
	return &UserAuth{ur: auth}
}

func (ua *UserAuth) CreateAuth(id uint64) (http.Cookie, error) {
	return ua.ur.CreateAuth(id)
}

func (ua *UserAuth) CheckAuth(sessionValue string) (uint64, error) {
	return ua.ur.CheckAuth(sessionValue)
}

func (ua *UserAuth) DeleteAuth(details *repository.AccessDetails) error {
	return ua.ur.DeleteAuth(details)
}
