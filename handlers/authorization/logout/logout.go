package logout

import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	session.MaxAge = -1
	http.SetCookie(w, session)
	w.WriteHeader(http.StatusFound)
}
