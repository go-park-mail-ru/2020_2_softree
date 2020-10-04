package authorization

import (
	"encoding/json"
	"net/http"
	"server/domain/entity"
	"server/infrastructure/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginJSON entity.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exist := UsersServerSession[loginJSON.Email]; !exist {
		var errorJSON entity.ErrorJSON

		errorJSON.Email = append(errorJSON.Email, "user does not exist")
		result, err := json.Marshal(errorJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result)
		http.Redirect(w, r, SignupPage, http.StatusBadRequest)

		return
	}

	if UsersServerSession[loginJSON.Email] != security.MakeDoubleHash(loginJSON.Password) {
		http.Redirect(w, r, LoginPage, http.StatusBadRequest)

		return
	}

	cookie := security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, RootPage, http.StatusFound)
}
