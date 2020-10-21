package persistence

import (
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

type UserRepo struct {
	database string
}

var users []entity.User

func NewUserRepository(database string) *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) SaveUser(u *entity.User) (*entity.User, *jsonRealisation.ErrorJSON) {
	users = append(users, *u)
	u.ID = uint64(len(users))
	return u, &jsonRealisation.ErrorJSON{}
}

func (ur *UserRepo) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	var user entity.User
	for _, user = range users {
		if user.ID == id {
			break
		}
	}

	user = u
	return user, nil
}

func (ur *UserRepo) DeleteUser(id uint64) error {
	return nil
}

func (ur *UserRepo) GetUser(id uint64) (*entity.User, error) {
	return &entity.User{}, nil
}
