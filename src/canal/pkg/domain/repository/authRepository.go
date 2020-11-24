package repository

import (
	"net/http"
)

type AuthRepository interface {
	authCreator
	authChecker
	authEraser
	authCookie
}

type authCreator interface {
	CreateAuth(int64) (http.Cookie, error)
}

type authChecker interface {
	CheckAuth(string) (int64, error)
}

type authEraser interface {
	DeleteAuth(string) error
}

type authCookie interface {
	CreateCookie() (http.Cookie, error)
}
