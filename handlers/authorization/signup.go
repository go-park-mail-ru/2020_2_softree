package authorization

import (
	"encoding/json"
	"net/http"
	"regexp"
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
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length,"+
		"Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func isValidEmail(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
		"[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(str)
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

	if !isValidEmail(signupJSON.Email) {
		var errorJSON entity.ErrorJSON

		errorJSON.Email = append(errorJSON.Email, "not an e-mail")
		result, err := json.Marshal(errorJSON)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result)
		http.Redirect(w, r, SignupPage, http.StatusBadRequest)

		return
	}

	if _, exist := UsersServerSession[signupJSON.Email]; exist {
		var errorJSON entity.ErrorJSON

		errorJSON.Email = append(errorJSON.Email, "user already exists")
		result, err := json.Marshal(errorJSON)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(result)
		http.Redirect(w, r, SignupPage, http.StatusBadRequest)

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

	cookie := security.MakeCookie(signupJSON.Email, r.Header.Get("Origin"))
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, RootPage, http.StatusOK)
}
