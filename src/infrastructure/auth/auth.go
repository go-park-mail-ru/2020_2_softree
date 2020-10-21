package auth

import (
	"net/http"
)

type AuthInterface interface {
	CreateAuth(uint64, *http.Cookie) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	SessionId string `json:"session_id"`
}
