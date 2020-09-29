package Security

import (
	"net/http"
	"time"
)

func MakeCookie(valueForCookie string) http.Cookie {
	expiration := time.Now().Add(10 * time.Hour)
	return http.Cookie{
		Name:     "session_id",
		Value:    MakeDoubleHash(valueForCookie),
		Expires:  expiration,
		HttpOnly: true,
	}
}
