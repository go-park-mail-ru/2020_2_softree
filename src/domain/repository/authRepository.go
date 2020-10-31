package repository

import "net/http"

type AuthRepository interface {
	creator
	check
	delete
}

type creator interface {
	CreateAuth(uint64) (http.Cookie, error)
}

type check interface {
	CheckAuth(string) (uint64, error)
}

type delete interface {
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	Value string `json:"session_id"`
}
