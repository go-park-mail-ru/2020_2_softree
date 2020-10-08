package signup

import (
	"encoding/json"
	"net/http"
	"server/src/domain/entity"
	"server/src/domain/entity/jsonRealisation"
	"server/src/handlers/authorization/utils"
	"server/src/infrastructure/security"
	"strings"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	signupJSON := new(jsonRealisation.SignupJSON)

	errorJSON := utils.Validate(signupJSON, w, r)
	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		utils.AddToErrorForm(&errorJSON, "NonFieldError", "Пароли не совпадают")
	}

	if errorJSON.NotEmpty { // contains some errors
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(errorJSON)

		_, err := w.Write(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	utils.UsersServerSession[signupJSON.Email] = security.MakeShieldedHash(signupJSON.Password1)
	cookie := security.MakeCookie()
	utils.Sessions[signupJSON.Email] = cookie.Value
	http.SetCookie(w, &cookie)

	entity.Users = append(entity.Users, entity.PublicUser{Email: signupJSON.Email})

	w.WriteHeader(http.StatusCreated)
}
