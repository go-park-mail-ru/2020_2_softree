package logout

import (
	"log"
	"net/http"
	"server/src/infrastructure/auth"
	"server/src/interfaces/authorization/utils"
	"server/src/interfaces/profile"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		w.WriteHeader(http.StatusFound)
		return
	}

	newCookie, err := auth.CreateCookie()
	if err != nil {
		log.Println(err)
		return
	}

	newCookie.Expires = time.Date(1973, 1, 1, 0, 0, 0, 0, time.UTC)
	newCookie.Value = ""
	http.SetCookie(w, &newCookie)

	email := profile.FindEmailInSession(cookie.Value)
	delete(utils.Sessions, email)

	w.WriteHeader(http.StatusFound)
}
