package auth

import "net/http"

type CookieInterface interface {
	CreateCookie(uint64) (*CookieDetails, error)
	ExtractData(*http.Request) (*AccessDetails, error)
}

type Cookie struct{}

func (c *Cookie) CreateCookie(userId uint64) (*CookieDetails, error) {
	return &CookieDetails{}, nil
}

func (c *Cookie) ExtractData(r *http.Request) (*AccessDetails, error) {
	return &AccessDetails{}, nil
}
