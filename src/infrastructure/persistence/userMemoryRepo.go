package persistence

import (
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
)

type UserMemoryRepo struct {
	users []entity.User
}

func NewUserRepository() *UserMemoryRepo {
	users := make([]entity.User, 1)
	return &UserMemoryRepo{users: users}
}

func (ur *UserMemoryRepo) SaveUser(u entity.User) (entity.User, error) {
	u.ID = uint64(len(ur.users) + 1)

	var err error
	u.Password, err = security.MakeShieldedPassword(u.Password)
	if err != nil {
		return entity.User{}, err
	}

	ur.users = append(ur.users, u)
	return u, nil
}

func (ur *UserMemoryRepo) UpdateUser(id uint64, u entity.User) (entity.User, error) {
	var user entity.User
	var i int
	for i, user = range ur.users {
		if user.ID == id {
			break
		}
	}

	if !govalidator.IsNull(u.Password) {
		ur.users[i].Password, _ = security.MakeShieldedPassword(u.Password)
		user.Password, _ = security.MakeShieldedPassword(u.Password)
	}
	if !govalidator.IsNull(u.Avatar) {
		ur.users[i].Avatar = u.Avatar
		user.Avatar = u.Avatar
	}

	return user, nil
}

func (ur *UserMemoryRepo) DeleteUser(id uint64) error {
	var user entity.User
	var i int
	for i, user = range ur.users {
		if user.ID == id {
			ur.users = append(ur.users[:i], ur.users[i + 1:]...)
		}
	}

	return nil
}

func (ur *UserMemoryRepo) GetUserById(id uint64) (entity.User, error) {
	var user entity.User
	for _, user = range ur.users {
		if user.ID == id {
			break
		}
	}

	return user, nil
}

func (ur *UserMemoryRepo) GetUserByLogin(email, password string) (entity.User, error) {
	var user entity.User
	for _, user = range ur.users {
		if user.Email == email && bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
			break
		}
	}

	return user, nil
}
