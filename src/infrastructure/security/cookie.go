package security

import (
	"net/http"
	"server/src/infrastructure/config"
	"time"
)

func MakeCookie() (http.Cookie, error) {
	hash, err := makeCookieHash()
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    hash,
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	}, nil
}
