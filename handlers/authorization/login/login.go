package login

import (
	"encoding/json"
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/auth"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if utils.UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.Sessions[loginJSON.Email] = security.MakeShieldedHash(loginJSON.Email)
	cookie := security.MakeCookie()
	http.SetCookie(w, &cookie)

	u := auth.FindUserInSession(cookie.Value)
	result, err := json.Marshal(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
