package authorization

import (
	"net/http"
	"server/handlers/authorization/utils"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, utils.RootPage, http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, utils.RootPage, http.StatusFound)
}
