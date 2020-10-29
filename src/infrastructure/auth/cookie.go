package auth

import (
	"encoding/json"
	"net/http"
	"server/src/infrastructure/config"
	"server/src/infrastructure/security"
	"time"
)

type TokenHandler interface {
	CreateCookie() (http.Cookie, error)
	ExtractData(*http.Request) (AccessDetails, error)
}

type Token struct {}

func NewToken() *Token {
	return &Token{}
}

func (t *Token) ExtractData(r *http.Request) (ad AccessDetails, err error) {
	if err := json.NewDecoder(r.Body).Decode(&ad); err != nil {
		return AccessDetails{}, err
	}
	return
}

func (t *Token) CreateCookie() (http.Cookie, error) {
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
