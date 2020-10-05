package signup

import (
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	signupJSON := new(jsonRealisation.SignupJSON)

	errorMas := utils.Validate(signupJSON, w, r)
	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		errorMas = append(errorMas, "Пароли не совпадают")
		utils.CreateErrorForm(w, errorMas)
	}

	if len(errorMas) != 0 {  // contain some errors
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	utils.Sessions[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Email)

	cookie := security.MakeCookie(signupJSON.Email)
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusCreated)
}
