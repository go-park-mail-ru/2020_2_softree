package auth

import (
	"errors"
	"net/http"
	"server/src/infrastructure/config"
	"server/src/infrastructure/security"
	"time"
)

type MemAuth struct {
	sessions []session
}

func NewMemAuth() *MemAuth {
	sessions := make([]session, 1)
	return &MemAuth{sessions: sessions}
}

type session struct {
	id    uint64
	value string
}

func (m *MemAuth) CreateAuth(id uint64) (cookie http.Cookie, err error) {
	if cookie, err = m.CreateCookie(); err != nil {
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

func (m *MemAuth) DeleteAuth(value string) error {
	for i, val := range m.sessions {
		if val.value == value {
			m.sessions = append(m.sessions[:i], m.sessions[i+1:]...)
		}
	}

	return nil
}

func (m *MemAuth) CreateCookie() (http.Cookie, error) {
	hash, err := security.MakeShieldedCookie()
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	}, nil
}
