package security

import (
	"net/http"
	"server/src/infrastructure/config"
	"time"
)

func MakeCookie() http.Cookie {
	return http.Cookie{
		Name:     "session_id",
		Value:    makeCookieHash(),
		Expires:  time.Now().Add(10 * 24 * time.Hour),
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   config.GlobalServerConfig.Secure,
		HttpOnly: true,
		Path:     "/",
	}
}
