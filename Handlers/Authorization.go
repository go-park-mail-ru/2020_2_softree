package Handlers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"server/Domain/Entity"
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

	var signupJSON Entity.SignupJSON
	if err := signupJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		http.Redirect(w, r, signup, http.StatusBadRequest)
		return
	}

	hash := md5.New()
	hash.Write([]byte(signupJSON.Password1))
	UsersServerSession[signupJSON.Email] = hex.EncodeToString(hash.Sum(nil))

	http.Redirect(w, r, login, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var loginJSON Entity.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := md5.New()
	hash.Write([]byte(loginJSON.Password))
	if UsersServerSession[loginJSON.Email] != hex.EncodeToString(hash.Sum(nil)) {
		http.Redirect(w, r, login, http.StatusBadRequest)
	}

	expiration := time.Now().Add(10 * time.Hour)
	hash.Write([]byte(loginJSON.Email))
	cookie := http.Cookie{
		Name:     "session_id",
		Value:    hex.EncodeToString(hash.Sum(nil)),
		Expires:  expiration,
		HttpOnly: true,
	}
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
