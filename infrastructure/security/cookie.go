package security

import (
	"net/http"
	"server/infrastructure/config"
	"time"
)

func MakeCookie(valueForCookie string) http.Cookie {
	expiration := time.Now().Add(10 * 24 * time.Hour)
	return http.Cookie{
		Name:     "session_id",
		Value:    MakeShieldedHash(valueForCookie),
		Expires:  expiration,
		Domain:   config.GlobalServerConfig.Domain,
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
	}
}
