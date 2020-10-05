package logout

import (
	"net/http"
	"server/handlers/authorization/utils"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", utils.RootPage)
		w.WriteHeader(http.StatusOK)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	w.Header().Set("Location", utils.RootPage)
	w.WriteHeader(http.StatusOK)
}
