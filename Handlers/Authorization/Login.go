package Authorization

import (
	"net/http"
	"server/Domain/Entity"
	"server/Infrastructure/Security"
)

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
		http.Redirect(w, r, LoginPage, http.StatusBadRequest)
	}

	cookie := Security.MakeCookie(loginJSON.Email)
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, RootPage, http.StatusFound)
}
