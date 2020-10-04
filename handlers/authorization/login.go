package authorization

import (
	"net/http"
	"server/domain/entity"
	"server/infrastructure/security"
)

func Login(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var loginJSON entity.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if UsersServerSession[loginJSON.Email] != security.MakeShieldedHash(loginJSON.Password) {
		http.Redirect(w, r, LoginPage, http.StatusBadRequest)
	}

	cookie := security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, RootPage, http.StatusFound)
}
