package logout

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	session.Expires = time.Date(1973, time.January, 1, 0, 0, 0, 0, time.UTC)
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusFound)
}
