package auth

import (
	"errors"
	"net/http"
	"server/src/domain/repository"
)

type MemAuth struct {
	token    TokenHandler
	sessions []session
}

func NewMemAuth() *MemAuth {
	sessions := make([]session, 1)
	return &MemAuth{sessions: sessions, token: NewToken()}
}

type session struct {
	id    uint64
	value string
}

func (m *MemAuth) CreateAuth(id uint64) (cookie http.Cookie, err error) {
	if cookie, err = m.token.CreateCookie(); err != nil {
		return http.Cookie{}, err
	}

	m.sessions = append(m.sessions, session{id: id, value: cookie.Value})
	return
}

func (m *MemAuth) CheckAuth(sessionValue string) (uint64, error) {
	for _, val := range m.sessions {
		if val.value == sessionValue {
			return val.id, nil
		}
	}

	return 0, errors.New("no session")
}

func (m *MemAuth) DeleteAuth(details *repository.AccessDetails) error {
	for i, val := range m.sessions {
		if val.value == details.Value {
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
		}
	}

	return nil
}
