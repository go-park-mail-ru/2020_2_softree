package security

import (
	"net/http"
	"time"
)

func MakeCookie(valueForCookie string, domain string) http.Cookie {
	expiration := time.Now().Add(10 * 24 * time.Hour)
	return http.Cookie{
		Name:     "session_id",
		Value:    MakeShieldedHash(valueForCookie),
		Expires:  expiration,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	}
}
