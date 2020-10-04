package authorization

import (
	"net/http"
	"server/domain/entity/jsonRealisation"
	"server/infrastructure/security"
	"strings"
)

const (
	LoginPage  = "/api/login"
	SignupPage = "/api/signup"
	RootPage   = "/"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	defer r.Body.Close()

	var signupJSON jsonRealisation.SignupJSON
	if validate(&signupJSON, &w, r) {
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		w.Header().Set("Location", SignupPage)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	Sessions[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Email)

	cookie := security.MakeCookie(signupJSON.Email, r.Header.Get("Origin"))
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, RootPage, http.StatusOK)
}
