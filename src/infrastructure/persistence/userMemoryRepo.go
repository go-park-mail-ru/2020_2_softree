package persistence

import (
	"github.com/asaskevich/govalidator"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
	"server/src/infrastructure/security"
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
	u.Password, _ = security.MakeShieldedHash(u.Password)
	Users = append(Users, u)
	return u, jsonRealisation.ErrorJSON{}
}

func (ur *UserMemoryRepo) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	var user entity.User
	var i int
	for i, user = range Users {
		if user.ID == id {
			break
		}
	}

	if !govalidator.IsNull(u.Password) {
		Users[i].Password, _ = security.MakeShieldedHash(u.Password)
		user.Password, _ = security.MakeShieldedHash(u.Password)
	}
	if !govalidator.IsNull(u.Avatar) {
		Users[i].Avatar = u.Avatar
		user.Avatar = u.Avatar
	}

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
