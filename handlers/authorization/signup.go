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
var Sessions = make(map[string]string, 0)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length," +
		"Accept-Encoding, X-CSRF-Token, Authorization")
}

func Signup(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == http.MethodOptions {
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

	if _, exist := UsersServerSession[signupJSON.Email]; exist {
		http.Redirect(w, r, SignupPage, http.StatusBadRequest)
		return
	}

	UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	Sessions[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Email)

	cookie := security.MakeCookie(signupJSON.Email)
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, RootPage, http.StatusOK)
}
