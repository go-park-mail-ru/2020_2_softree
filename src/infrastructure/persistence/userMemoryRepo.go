package persistence

import (
	"github.com/asaskevich/govalidator"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
)

type UserMemoryRepo struct {
	Users []entity.User
}

func NewUserRepository() *UserMemoryRepo {
	users := make([]entity.User, 1)
	return &UserMemoryRepo{Users: users}
}

func (ur *UserMemoryRepo) SaveUser(u entity.User) (entity.User, error) {
	u.ID = uint64(len(ur.Users) + 1)

	var err error
	u.Password, err = security.MakeShieldedHash(u.Password)
	if err != nil {
		return entity.User{}, err
	}

	ur.Users = append(ur.Users, u)
	return u, nil
}

func (ur *UserMemoryRepo) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	var user entity.User
	var i int
	for i, user = range ur.Users {
		if user.ID == id {
			break
		}
	}

	if !govalidator.IsNull(u.Password) {
		ur.Users[i].Password, _ = security.MakeShieldedHash(u.Password)
		user.Password, _ = security.MakeShieldedHash(u.Password)
	}
	if !govalidator.IsNull(u.Avatar) {
		ur.Users[i].Avatar = u.Avatar
		user.Avatar = u.Avatar
	}

	return user, nil
}

func (ur *UserMemoryRepo) DeleteUser(id uint64) error {
	var user entity.User
	var i int
	for i, user = range ur.Users {
		if user.ID == id {
			ur.Users = append(ur.Users[:i], ur.Users[i + 1:]...)
		}
	}

	return nil
}

func (ur *UserMemoryRepo) GetUser(id uint64) (entity.User, error) {
	var user entity.User
	for _, user = range ur.Users {
		if user.ID == id {
			break
		}
	}

	return user, nil
}
