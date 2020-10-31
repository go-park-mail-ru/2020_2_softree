package repository

import "net/http"

type AuthRepository interface {
	AuthCreator
	AuthChecker
	AuthEraser
}

type AuthCreator interface {
	CreateAuth(uint64) (http.Cookie, error)
}

type AuthChecker interface {
	CheckAuth(string) (uint64, error)
}

type AuthEraser interface {
	DeleteAuth(*AccessDetails) error
}

type AccessDetails struct {
	Value string `json:"session_id"`
}
