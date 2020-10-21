package auth

import (
	"errors"
)

type AuthInterface interface {
	CreateAuth(uint64, string) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	SessionId uint64 `json:"session_id"`
}

type MemAuth struct {
	auth string
}

func NewMemAuth(auth string) *MemAuth {
	return &MemAuth{auth: auth}
}

type Session struct {
	ID    uint64
	Value string
}

var Sessions []Session

func (m *MemAuth) CreateAuth(id uint64, sessionValue string) error {
	Sessions = append(Sessions, Session{ID: id, Value: sessionValue})
	return nil
}

func (m *MemAuth) CheckAuth(sessionValue string) (uint64, error) {
	for _, val := range Sessions {
		if val.Value == sessionValue {
			return val.ID, nil
		}
	}

	return 0, errors.New("no session")
}

func (m *MemAuth) DeleteAuth(details *AccessDetails) error {
	for i, val := range Sessions {
		if val.ID == details.SessionId {
			Sessions = append(Sessions[:i], Sessions[i + 1:]...)
		}
	}
	return nil
}
