package login

import (
	"encoding/json"
	"log"
	"net/http"
	"server/src/domain/entity/jsonRealisation"
	authInf "server/src/infrastructure/auth"
	"server/src/infrastructure/security"
	authAuth"server/src/interfaces/authorization/auth"
	"server/src/interfaces/authorization/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var loginJSON jsonRealisation.LoginJSON
	if err := loginJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exists := utils.UsersServerSession[loginJSON.Email]; !exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userTryToLogin, _ := security.MakeShieldedHash(loginJSON.Password)
	if utils.UsersServerSession[loginJSON.Email] != userTryToLogin {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookie, err := authInf.CreateCookie()
	if err != nil {
		log.Println(err)
		return
	}
	utils.Sessions[loginJSON.Email] = cookie.Value
	http.SetCookie(w, &cookie)

	u, _ := authAuth.FindUserInSession(cookie.Value)
	result, err := json.Marshal(&u)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(result)
	if err != nil {
		log.Println(err)
		return
	}
}
