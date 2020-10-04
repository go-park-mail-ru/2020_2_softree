package authorization

import (
	"net/http"
	"server/domain/entity"
	"server/infrastructure/security"
	"strings"
)

const (
	LoginPage  = "/api/login"
	SignupPage = "/api/signup"
	RootPage   = "/"
)

var UsersServerSession = make(map[string]string, 0)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var signupJSON entity.SignupJSON
	if err := signupJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		http.Redirect(w, r, SignupPage, http.StatusBadRequest)
		return
	}

	UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	cookie := security.MakeCookie(signupJSON.Email)
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, RootPage, http.StatusOK)
}
