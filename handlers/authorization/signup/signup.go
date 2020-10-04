package signup

import (
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	utils.EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	signupJSON := new(jsonRealisation.SignupJSON)
	if !utils.Validate(signupJSON, &w, r) {
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		w.Header().Set("Location", utils.SignupPage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	utils.UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	utils.Sessions[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Email)

	cookie := security.MakeCookie(signupJSON.Email, r.Header.Get("Origin"))

	http.SetCookie(w, &cookie)
	w.Header().Set("Location", utils.RootPage)
	w.WriteHeader(http.StatusOK)
}
