package userInteraction

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
	"server/src/handlers/authorization/utils"
	"server/src/infrastructure/security"
)

func UpdateUserPartly(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var userJSON jsonRealisation.UserJSON
		err := userJSON.FillFields(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		changeAvatar(*cookie, userJSON)

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	if logged {
		var userJSON jsonRealisation.UserJSON
		if err := userJSON.FillFields(r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		code, errorJSON := changePassword(*cookie, userJSON)
		if code == http.StatusOK {
			w.WriteHeader(code)
			return
		}

		res, _ := json.Marshal(errorJSON)
		w.WriteHeader(code)
		if _, err := w.Write(res); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func FindEmailInSession(hash string) string {
	for key, val := range utils.Sessions {
		if val == hash {
			return key
		}
	}

	return ""
}

func changePassword(cookie http.Cookie, userJSON jsonRealisation.UserJSON) (int, jsonRealisation.ErrorJSON) {
	emailInSession := FindEmailInSession(cookie.Value)
	userPassword := utils.UsersServerSession[emailInSession]
	oldPassword, _ := security.MakeShieldedHash(userJSON.OldPassword)

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
		return http.StatusBadRequest, errorJSON
	}

	utils.UsersServerSession[emailInSession], _ = security.MakeShieldedHash(userJSON.NewPassword1)
	return http.StatusOK, errorJSON
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
