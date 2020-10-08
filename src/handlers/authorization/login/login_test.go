package login

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/src/domain/entity/jsonRealisation"
	"server/src/handlers/authorization/utils"
	"server/src/infrastructure/security"
	"strings"
	"testing"
)

func TestLoginFailWithGet(t *testing.T) {
	url := "http://example.com/api/"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusBadRequest)
	}
}

func TestLoginSuccess(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	utils.UsersServerSession["yandex@mail.ru"], _ = security.MakeShieldedHash("str")

	Login(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Errorf("no cookie")
	}

	if len(utils.UsersServerSession) == 0 {
		t.Errorf("no users created")
	}

	if len(utils.Sessions) == 0 {
		t.Errorf("no session")
	}

	delete(utils.UsersServerSession, "yandex@mail.ru")
	delete(utils.Sessions, "yandex@mail.ru")
}

func TestLoginUserNonExist(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("unexpected cookie was created")
	}

	if len(utils.Sessions) != 0 {
		t.Errorf("unexpected session was created")
	}
}

func TestSignupFailPassword(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	utils.UsersServerSession["yandex@mail.ru"], _ = security.MakeShieldedHash("string")

	Login(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("unexpected cookie was created")
	}

	if len(utils.Sessions) != 0 {
		t.Errorf("unexpected session was created")
	}

	delete(utils.UsersServerSession, "yandex@mail.ru")
	delete(utils.Sessions, "yandex@mail.ru")
}

func TestLoginFailWithNotFilledField(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := jsonRealisation.SignupJSON{
		Email:     "right",
		Password1: "str",
		Password2: "ste",
	}
	body, _ := json.Marshal(jsonForBody)
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()

	Login(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusBadRequest)
	}
}
