package signup

import (
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	signupJSON := new(jsonRealisation.SignupJSON)
	if !utils.Validate(signupJSON, &w, r) {
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		errorMas := make([]string, 0)
		errorMas = append(errorMas, "fail to compare passwords")
		utils.CreateErrorForm(&w, errorMas)
		return
	}

	utils.UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	utils.Sessions[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Email)

	cookie := security.MakeCookie(signupJSON.Email, r.Header.Get("Origin"))

	http.SetCookie(w, &cookie)
	w.Header().Set("Location", utils.RootPage)
	w.WriteHeader(http.StatusOK)
}
