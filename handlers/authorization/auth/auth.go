package auth

import (
	"encoding/json"
	"net/http"
	"server/domain/entity"
	"server/handlers/authorization/utils"
)

func Authentication(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		u := FindUserInSession(cookie.Value)

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

func FindUserInSession(hash string) entity.PublicUser {
	var email string
	for key, val := range utils.Sessions {
		if val == hash {
			email = key
		}
	}

	for i, _ := range entity.Users {
		if entity.Users[i].Email == email {
			return entity.Users[i]
		}
	}

	return entity.PublicUser{}
}
