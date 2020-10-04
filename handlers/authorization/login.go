package authorization

import (
	"encoding/json"
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
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
			http.Redirect(w, r, utils.SignupPage, http.StatusBadRequest)

			return
		}

		if utils.UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {
			http.Redirect(w, r, utils.LoginPage, http.StatusBadRequest)

			return
		}

		cookie := security.MakeCookie(loginJSON.Email, r.Header.Get("Origin"))
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, utils.RootPage, http.StatusFound)
	}
}
