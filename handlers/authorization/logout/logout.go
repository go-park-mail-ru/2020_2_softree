package logout

import (
	"net/http"
	"server/handlers/authorization/utils"
	"server/handlers/userInteraction"
	"server/infrastructure/security"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	newCookie := security.MakeCookie()
	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	http.SetCookie(w, &newCookie)

	email := userInteraction.FindEmailInSession(cookie.Value)
	delete(utils.Sessions, email)

	w.WriteHeader(http.StatusFound)
}
