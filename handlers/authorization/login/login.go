package login

import (
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	var loginJSON jsonRealisation.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exists := utils.UsersServerSession[loginJSON.Email]; !exists {
		errorMas := []string {"Пользователь не существует"}
		utils.CreateErrorForm(w, errorMas)
		return
	}

	if utils.UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {
		errorMas := []string {"incorrect password"}
		utils.CreateErrorForm(w, errorMas)
		return
	}

	utils.Sessions[loginJSON.Email] = security.MakeShieldedHash(loginJSON.Email)
	cookie := security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
