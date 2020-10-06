package userInteraction

import (
	"net/http"
	"server/domain/entity"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var userJSON jsonRealisation.UserJSON
		userJSON.FillFields(r.Body)

		if r.Method == "PUT" {
			changeEmail(*cookie, userJSON)
			changeAvatar(*cookie, userJSON)
		} else if r.Method == "PATCH" {
			changeAvatar(*cookie, userJSON)
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func findEmailInSession(hash string) string {
	for key, val := range utils.Sessions {
		if val == hash {
			return key
		}
	}

	return ""
}

func changeEmail(cookie http.Cookie, userJSON jsonRealisation.UserJSON) {
	emailInSession := findEmailInSession(cookie.Value)
	userHashedPassword := utils.UsersServerSession[emailInSession]

	delete(utils.Sessions, emailInSession)
	utils.Sessions[userJSON.Email] = cookie.Value

	delete(utils.UsersServerSession, emailInSession)
	utils.UsersServerSession[userJSON.Email] = userHashedPassword

	for i, _ := range entity.Users {
		if entity.Users[i].Email == emailInSession {
			entity.Users[i].Email = userJSON.Email
			break
		}
	}
}

func changeAvatar(cookie http.Cookie, userJSON jsonRealisation.UserJSON) {
	emailInSession := findEmailInSession(cookie.Value)
	for i, _ := range entity.Users {
		if entity.Users[i].Email == emailInSession {
			entity.Users[i].Avatar = userJSON.Avatar
			break
		}
	}
}
