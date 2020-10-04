package UserInteraction

import (
	"encoding/json"
	"net/http"
	"server/Domain/Entity"
	"server/Handlers/Authorization"
	"strconv"
)

const (
	testID    = 123
	testEmail = "hound@psina.ru"
)

func UserData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := r.Cookie("session_id")
	logged := err != http.ErrNoCookie

	var u Entity.User
	u.ID = testID
	u.Email = testEmail
	result, e := json.Marshal(u)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if logged {
		w.Write(result)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, Authorization.SignupPage, http.StatusUnauthorized)
	}
}

func Rates(w http.ResponseWriter, r *http.Request) {
	result := make([]byte, 0)
	id, err := strconv.Atoi(r.FormValue("id"))
	if id == 0 {
		var quotations Entity.Quotations
		quotations.Add(
			Entity.CurrencyQuotation{Value: 1234, Title: "title1", Change: 12},
			Entity.CurrencyQuotation{Value: 12, Title: "title2", Change: 1},
		)

		result, err = json.Marshal(quotations)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		result, err = json.Marshal(Entity.FindById(uint64(id)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Write(result)
	w.WriteHeader(http.StatusOK)
}
