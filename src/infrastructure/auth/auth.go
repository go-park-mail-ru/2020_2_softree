package auth

import (
	"net/http"
	"time"
)

type AuthInterface interface {
	CreateAuth(uint64, *CookieDetails) error
	CheckAuth(string) (uint64, error)
	DeleteAuth(*AccessDetails) error
}

type CookieDetails struct {
	Name     string
	Value    string
	Expires  time.Time
	Domain   string
	Secure   bool
	HttpOnly bool
	Path     string
}

type AccessDetails struct {
	SessionId string `json:"session_id"`
}

func MakeCookieDetailsFromCookie(cookie http.Cookie) CookieDetails {
	return CookieDetails{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Expires:  cookie.Expires,
		Domain:   cookie.Domain,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
		Path:     cookie.Path,
	}
}
