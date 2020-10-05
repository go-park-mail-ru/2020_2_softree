package signup

import (
	"encoding/json"
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	signupJSON := new(jsonRealisation.SignupJSON)

	errorJSON := utils.Validate(signupJSON, w, r)
	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		utils.AddToErrorForm(&errorJSON, "NonFieldError", "Пароли не совпадают")
	}

	if errorJSON.NotEmpty {  // contains some errors
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(errorJSON)
		w.Write(res)
		return
	}

	utils.UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	cookie := security.MakeCookie(signupJSON.Email)
	utils.Sessions[signupJSON.Email] = cookie.Value
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}
