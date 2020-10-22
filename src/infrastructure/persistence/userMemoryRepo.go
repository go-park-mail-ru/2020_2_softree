package persistence

import (
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

type UserMemoryRepo struct {
	database string
}

var Users []entity.User

func NewUserRepository(database string) *UserMemoryRepo {
	return &UserMemoryRepo{database: database}
}

func (ur *UserMemoryRepo) SaveUser(u entity.User) (entity.User, jsonRealisation.ErrorJSON) {
	u.ID = uint64(len(Users) + 1)
	Users = append(Users, u)
	return u, jsonRealisation.ErrorJSON{}
}

func (ur *UserMemoryRepo) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	var user entity.User
	for _, user = range Users {
		if user.ID == id {
			break
		}
	}

	user = u
	return user, nil
}

func (ur *UserMemoryRepo) DeleteUser(id uint64) error {
	return nil
}

func (ur *UserMemoryRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	for _, user = range Users {
		if user.ID == id {
			break
		}
	}

	return &user, nil
}
