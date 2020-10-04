package authorization

import (
	"encoding/json"
	"net/http"
	"regexp"
	"server/domain/entity"
)

var UsersServerSession = make(map[string]string, 0)
var Sessions = make(map[string]string, 0)

func validate(JSON *entity.JSON, w *http.ResponseWriter, r *http.Request) {
	if err := (*JSON).FillFields(r.Body); err != nil {
		(*w).WriteHeader(http.StatusBadRequest)
		return
	}

	errorMas := make([]string, 0)
	if !isValidEmail((*JSON).GetEmail()) {
		errorMas = append(errorMas, "not an e-mail")
		createErrorForm(&(*w), errorMas)
		return
	}
	if _, exist := UsersServerSession[(*JSON).GetEmail()]; exist {
		errorMas = append(errorMas, "user already exists")
		createErrorForm(&(*w), errorMas)
		return
	}
}

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

func createErrorForm(w *http.ResponseWriter, messages []string) {
	var errorJSON entity.ErrorJSON

	errorJSON.Email = append(errorJSON.Email, messages...)
	result, err := json.Marshal(errorJSON)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return
	}
	(*w).Header().Set("Location", SignupPage)
	(*w).WriteHeader(http.StatusBadRequest)
	(*w).Write(result)
}
