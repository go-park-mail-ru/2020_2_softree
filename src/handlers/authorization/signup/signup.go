package signup

import (
	"encoding/json"
	"log"
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
	if err := signupJSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	errorJSON := utils.Validate(signupJSON)
	if strings.Compare(signupJSON.Password1, signupJSON.Password2) != 0 {
		utils.AddToErrorForm(&errorJSON, "NonFieldError", "Пароли не совпадают")
	}

	if errorJSON.NotEmpty { // contains some errors
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(errorJSON)

		_, err := w.Write(res)
		if err != nil {
			log.Println(err)
			return
		}

		return
	}

	var err error
	utils.UsersServerSession[signupJSON.Email], err = security.MakeShieldedHash(signupJSON.Password1)
	if err != nil {
		log.Println(err)
		return
	}

	cookie, err := security.MakeCookie()
	if err != nil {
		log.Println(err)
		return
	}

	utils.Sessions[signupJSON.Email] = cookie.Value
	http.SetCookie(w, &cookie)
	entity.Users = append(entity.Users, entity.PublicUser{Email: signupJSON.Email})

	w.WriteHeader(http.StatusCreated)
}
