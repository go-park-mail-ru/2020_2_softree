package Handlers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

const (
	login  = "/api/login"
	signup = "/api/signup"
	root   = "/"
)

var UsersServerSession map[string]string

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password1 := r.FormValue("password1")
	password2 := r.FormValue("password2")

	if strings.Compare(password1, password2) != 0 {
		http.Redirect(w, r, signup, http.StatusBadRequest)
	}

	hash := md5.New()
	hash.Write([]byte(password1))
	UsersServerSession[email] = hex.EncodeToString(hash.Sum(nil))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Location", login)
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash := md5.New()
	hash.Write([]byte(password))
	if UsersServerSession[email] != hex.EncodeToString(hash.Sum(nil)) {
		http.Redirect(w, r, login, http.StatusBadRequest)
	}

	expiration := time.Now().Add(10 * time.Hour)
	hash.Write([]byte(email))
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
