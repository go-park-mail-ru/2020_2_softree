package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/domain/entity"
	"server/handlers/authorization/utils"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var u entity.PublicUser
		u.Email = findEmailInSession(cookie.Value)

		result, err := json.Marshal(&u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(result)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func findEmailInSession(hash string) string {
	for key, val := range utils.Sessions {
		if val == hash {
			fmt.Println(key)
			return key
		}
	}

	return ""
}
