package Handlers

import (
	"net/http"
	"server/Domain/Entity"
	"server/Infrastructure/Security"
	"strings"
	"time"
)

const (
	login  = "/api/login"
	signup = "/api/signup"
	root   = "/"
)

var UsersServerSession = make(map[string]string, 0)

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var signupJSON Entity.SignupJSON
	if err := signupJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		http.Redirect(w, r, signup, http.StatusBadRequest)
		return
	}

	UsersServerSession[signupJSON.Email] = Security.MakeDoubleHash(signupJSON.Password1)
	http.Redirect(w, r, login, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginJSON Entity.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if UsersServerSession[loginJSON.Email] != Security.MakeDoubleHash(loginJSON.Password) {
		http.Redirect(w, r, login, http.StatusBadRequest)
	}

	cookie := Security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, root, http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		http.Redirect(w, r, root, http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)

	http.Redirect(w, r, root, http.StatusFound)
}
