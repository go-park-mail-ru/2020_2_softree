package utils

import (
	"encoding/json"
	"net/http"
	"regexp"
	"server/domain/entity/jsonRealisation"
)

func Validate(JSON jsonRealisation.JSON, w http.ResponseWriter, r *http.Request) []string {
	if err := JSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	errorMas := make([]string, 0)
	if !isValidEmail(JSON.GetEmail()) {
		errorMas = append(errorMas, "Неправильный формат Email")
		CreateErrorForm(w, errorMas)
	}
	if _, exist := UsersServerSession[JSON.GetEmail()]; exist {
		errorMas = append(errorMas, "Пользователь с таким Email уже существует")
		CreateErrorForm(w, errorMas)
	}

	return errorMas
}

func isValidEmail(str string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
		"[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(str)
}

func CreateErrorForm(w http.ResponseWriter, messages []string) {
	var errorJSON jsonRealisation.ErrorJSON

	errorJSON.Email = append(errorJSON.Email, messages...)
	result, err := json.Marshal(errorJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(result)
}
