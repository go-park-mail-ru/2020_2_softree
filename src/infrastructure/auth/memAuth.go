package auth

import (
	"errors"
	"server/src/domain/repository"
)

type MemAuth struct {
	Sessions []Session
}

func NewMemAuth() *MemAuth {
	sessions := make([]Session, 1)
	return &MemAuth{Sessions: sessions}
}

type Session struct {
	ID    uint64
	Value string
}

func (m *MemAuth) CreateAuth(id uint64, sessionValue string) error {
	m.Sessions = append(m.Sessions, Session{ID: id, Value: sessionValue})
	return nil
}

func (m *MemAuth) CheckAuth(sessionValue string) (uint64, error) {
	for _, val := range m.Sessions {
		if val.Value == sessionValue {
			return val.ID, nil
		}
	}

	return 0, errors.New("no session")
}

func (m *MemAuth) DeleteAuth(details *repository.AccessDetails) error {
	for i, val := range m.Sessions {
		if val.Value == details.Value {
			m.Sessions = append(m.Sessions[:i], m.Sessions[i+1:]...)
		}
	}
	return nil
}
