package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"server/src/domain/entity"
	"server/src/handlers/authorization/utils"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		u := FindUserInSession(cookie.Value)

		result, _ := json.Marshal(&u)

		w.Header().Set("content-type", "application/json")

		_, err := w.Write(result)
		if err != nil {
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func FindUserInSession(hash string) entity.PublicUser {
	var email string
	for key, val := range utils.Sessions {
		if val == hash {
			email = key
		}
	}

	for i := range entity.Users {
		if entity.Users[i].Email == email {
			return entity.Users[i]
		}
	}

	return entity.PublicUser{}
}
