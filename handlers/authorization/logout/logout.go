package logout

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
<<<<<<< HEAD
		w.Header().Set("Location", utils.RootPage)
		w.WriteHeader(http.StatusOK)
=======
		w.WriteHeader(http.StatusFound)
>>>>>>> origin/michael
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
<<<<<<< HEAD

	w.Header().Set("Location", utils.RootPage)
	w.WriteHeader(http.StatusOK)
=======
	w.WriteHeader(http.StatusFound)
>>>>>>> origin/michael
}
