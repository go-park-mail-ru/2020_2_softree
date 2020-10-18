package persistence

import (
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
)

type UserRepo struct {
	database string  // *gorm.DB
}

func NewUserRepository(database string) *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) SaveUser(u *entity.User) (*entity.User, *jsonRealisation.ErrorJSON) {
	return u, &jsonRealisation.ErrorJSON{}
}

func (ur *UserRepo) UpdateUser(uint64, entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (ur *UserRepo) DeleteUser(uint64) error {
	return nil
}

func (ur *UserRepo) GetUser(uint64) (*entity.User, error) {
	return &entity.User{}, nil
}
