package logout

import (
	"net/http"
	"server/handlers/authorization/utils"
	"server/handlers/userInteraction"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	cookie.Expires = time.Date(1973, time.January, 1, 0, 0, 0, 0, time.UTC)
	cookie.Value = ""
	cookie.MaxAge = 0
	http.SetCookie(w, cookie)

	email := userInteraction.FindEmailInSession(cookie.Value)
	delete(utils.Sessions, email)

	w.WriteHeader(http.StatusFound)
}
