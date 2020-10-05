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

	if _, isRegistered := utils.UsersServerSession[loginJSON.Email]; !isRegistered {
		errorMas := []string {"user does not exist"}
		utils.CreateErrorForm(&w, errorMas)
		return
	}

	if utils.UsersServerSession[loginJSON.Password] != security.MakeShieldedHash(loginJSON.Password) {
		errorMas := []string {"incorrect password"}
		utils.CreateErrorForm(&w, errorMas)
		return
	}

	cookie := security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
