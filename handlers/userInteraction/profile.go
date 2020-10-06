package userInteraction

import (
	"encoding/json"
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
			changePassword(*cookie, userJSON, w)
			changeAvatar(*cookie, userJSON)
		} else {
			if r.Method == "PATCH" {
				if userJSON.Avatar != "" {
					// change only avatar
					changeAvatar(*cookie, userJSON)
				} else {
					// change only password
					changePassword(*cookie, userJSON, w)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var userJSON jsonRealisation.UserJSON
		userJSON.FillFields(r.Body)

		if changePassword(*cookie, userJSON, w) {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
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

func changePassword(cookie http.Cookie, userJSON jsonRealisation.UserJSON, w http.ResponseWriter) bool {
	emailInSession := findEmailInSession(cookie.Value)
	userPassword := utils.UsersServerSession[emailInSession]

	var errorJSON jsonRealisation.ErrorJSON
	if userPassword != userJSON.OldPassword {
		errorJSON.NotEmpty = true
		errorJSON.OldPassword = append(errorJSON.OldPassword, "Введен неверно старый пароль")
	}
	if userJSON.NewPassword1 != userJSON.NewPassword2 {
		errorJSON.NotEmpty = true
		errorJSON.NonFieldError = append(errorJSON.NonFieldError, "Пароли не совпадают")
	}

	if errorJSON.NotEmpty {
		res, _ := json.Marshal(errorJSON)
		w.Write(res)
		return false
	}

	utils.UsersServerSession[emailInSession] = userJSON.NewPassword1

	return true
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
