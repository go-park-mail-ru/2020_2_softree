package auth

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"server/domain/entity"
	"server/handlers/authorization/utils"
	"server/infrastructure/security"
	"strings"
	"testing"
)

func TestAuthenticationFail(t *testing.T) {
	url := "http://example.com/api/"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Authentication(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticationFailMarshal(t *testing.T) {
	url := "http://example.com/api/"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	cookie := security.MakeCookie()
	req.AddCookie(&cookie)
	utils.Sessions["yandex@mail.ru"] = cookie.Value

	Authentication(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusOK)
	}
}

func TestAuthenticationSuccess(t *testing.T) {
	url := "http://example.com/api/"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	cookie := security.MakeCookie()
	req.AddCookie(&cookie)
	utils.Sessions["yandex@mail.ru"] = cookie.Value

	Authentication(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusOK)
	}

	expected := "{\"email\":\"\",\"avatar\":\"\"}"

	bodyBytes, _ := ioutil.ReadAll(w.Result().Body)
	bodyString := string(bodyBytes)

	if bodyString != expected {
		t.Errorf("\nwrong response body\ngot: %s\nexpected: %s",
			bodyString, expected)
	}
}

func TestFindUserInSessionSuccess(t *testing.T) {
	url := "http://example.com/api/"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	cookie := security.MakeCookie()
	req.AddCookie(&cookie)

	val := cookie.Value
	utils.Sessions["yandex@mail.ru"] = val
	user := entity.PublicUser{Email: "yandex@mail.ru", Avatar: "some"}
	entity.Users = append(entity.Users, user)

	result := FindUserInSession(val)

	if result != user {
		t.Errorf("\nwrong result\ngot: %s\nexpected: %s",
			result, user)
	}
}

func TestFindUserInSessionFail(t *testing.T) {
	url := "http://example.com/api/"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	cookie := security.MakeCookie()
	req.AddCookie(&cookie)

	val := cookie.Value

	result := FindUserInSession(val)
	expected := entity.PublicUser{}

	if result != expected {
		t.Errorf("\nwrong result\ngot: %s\nexpected: %s",
			result, expected)
	}
}
