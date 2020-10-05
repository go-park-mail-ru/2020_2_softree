package login

import (
	"net/http"
	"net/http/httptest"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
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
	utils.UsersServerSession["yandex@mail.ru"] = security.MakeShieldedHash("str")

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
		t.Errorf("no cookie")
	}

	if len(utils.UsersServerSession) != 0 {
		t.Errorf("no users created")
	}

	if len(utils.Sessions) != 0 {
		t.Errorf("no session")
	}
}

func TestSignupFailPassword(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	utils.UsersServerSession["yandex@mail.ru"] = security.MakeShieldedHash("string")

	Login(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	cookies := w.Result().Cookies()
	if len(cookies) != 0 {
		t.Errorf("no cookie")
	}

	if len(utils.UsersServerSession) != 0 {
		t.Errorf("no users created")
	}

	if len(utils.Sessions) == 0 {
		t.Errorf("no session")
	}

	delete(utils.UsersServerSession, "yandex@mail.ru")
	delete(utils.Sessions, "yandex@mail.ru")
}
