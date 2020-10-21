package persistence

import (
	"server/src/domain/repository"
)

type Repositories struct {
	UserRepository repository.UserRepository
}

func NewRepositories() (*Repositories, error) {
	return &Repositories{}, nil
}

func (r *Repositories) Close() error {
	return nil
}

func (r *Repositories) Migrate() error {
	return nil
}
