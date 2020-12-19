package authorization

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func CreateCookie() http.Cookie {
	return http.Cookie{
		Name:     "session_id",
		Expires:  time.Now().Add(24 * time.Hour),
		Domain:   viper.GetString("server.domain"),
		Secure:   viper.GetBool("server.secure"),
		HttpOnly: true,
		Path:     "/",
	}
}
