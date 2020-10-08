package userInteraction

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
	"server/src/handlers/authorization/utils"
	"server/src/infrastructure/security"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var userJSON jsonRealisation.UserJSON
		userJSON.FillFields(r.Body)
		if r.Method == "PUT" {
			changeAvatar(*cookie, userJSON)
		} else if r.Method == "PATCH" {
			changeAvatar(*cookie, userJSON)
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

func FindEmailInSession(hash string) string {
	for key, val := range utils.Sessions {
		if val == hash {
			return key
		}
	}

	return ""
}

func changePassword(cookie http.Cookie, userJSON jsonRealisation.UserJSON, w http.ResponseWriter) bool {
	emailInSession := FindEmailInSession(cookie.Value)
	userPassword := utils.UsersServerSession[emailInSession]
	oldPassword := security.MakeShieldedHash(userJSON.OldPassword)

	var errorJSON jsonRealisation.ErrorJSON
	if userPassword != oldPassword {
		errorJSON.NotEmpty = true
		errorJSON.OldPassword = append(errorJSON.OldPassword, "Введен неверно старый пароль")
	}
	if userJSON.NewPassword1 != userJSON.NewPassword2 {
		errorJSON.NotEmpty = true
		errorJSON.NonFieldError = append(errorJSON.NonFieldError, "Пароли не совпадают")
	}

	if errorJSON.NotEmpty {
		res, _ := json.Marshal(errorJSON)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(res)
		return false
	}

	utils.UsersServerSession[emailInSession] = security.MakeShieldedHash(userJSON.NewPassword1)

	return true
}

func changeAvatar(cookie http.Cookie, userJSON jsonRealisation.UserJSON) {
	emailInSession := FindEmailInSession(cookie.Value)
	for i := range entity.Users {
		if entity.Users[i].Email == emailInSession {
			entity.Users[i].Avatar = userJSON.Avatar
			break
		}
	}
}
