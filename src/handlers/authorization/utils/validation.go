package utils

import (
	"net/http"
	"regexp"
	"server/src/domain/entity/jsonRealisation"
)

func Validate(JSON jsonRealisation.JSON, w http.ResponseWriter, r *http.Request) jsonRealisation.ErrorJSON {
	if err := JSON.FillFields(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	errorMas := make([]string, 0)
	if !isValidEmail(JSON.GetEmail()) {
		errorMas = append(errorMas, "Неправильный формат Email")
	}
	if _, exist := UsersServerSession[JSON.GetEmail()]; exist {
		errorMas = append(errorMas, "Пользователь с таким Email уже существует")
	}

	return CreateErrorForm("Email", errorMas...)
}

func isValidEmail(str string) bool {
	regStr := "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$"
	re := regexp.MustCompile(regStr)

	return re.MatchString(str)
}

func CreateErrorForm(errorType string, messages ...string) jsonRealisation.ErrorJSON {
	var errorJSON jsonRealisation.ErrorJSON
	if len(messages) == 0 { // no errors
		return errorJSON
	}

	errorJSON.NotEmpty = true

	switch errorType {
	case "Name":
		errorJSON.Name = append(errorJSON.Name, messages...)
	case "Email":
		errorJSON.Email = append(errorJSON.Email, messages...)
	case "Password":
		errorJSON.Password = append(errorJSON.Password, messages...)
	case "NonFieldError":
		errorJSON.NonFieldError = append(errorJSON.NonFieldError, messages...)
	}

	return errorJSON
}

func AddToErrorForm(errorJSON *jsonRealisation.ErrorJSON, errorType string, messages ...string) {
	errorJSON.NotEmpty = true

	switch errorType {
	case "Name":
		errorJSON.Name = append(errorJSON.Name, messages...)
	case "Email":
		errorJSON.Email = append(errorJSON.Email, messages...)
	case "Password":
		errorJSON.Password = append(errorJSON.Password, messages...)
	case "NonFieldError":
		errorJSON.NonFieldError = append(errorJSON.NonFieldError, messages...)
	}
}
