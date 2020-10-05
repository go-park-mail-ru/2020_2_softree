package login

import (
	"encoding/json"
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginJSON jsonRealisation.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if utils.UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {
		if _, exist := utils.UsersServerSession[loginJSON.Email]; !exist {
			var errorJSON jsonRealisation.ErrorJSON

			errorJSON.Email = append(errorJSON.Email, "user does not exist")
			result, err := json.Marshal(errorJSON)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(result)

			return
		}

		if utils.UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {

			return
		}

		cookie := security.MakeCookie(loginJSON.Email)
		http.SetCookie(w, &cookie)
	}
}
