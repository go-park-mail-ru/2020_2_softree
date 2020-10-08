package logout

import (
	"net/http"
	"net/http/httptest"
	"server/src/infrastructure/security"
	"strings"
	"testing"
	"time"
)

func TestLogoutFail(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Logout(w, req)

	if w.Result().StatusCode != http.StatusFound {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusFound)
	}
}

func TestLogoutSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"

	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	cookie := security.MakeCookie()
	req.AddCookie(&cookie)

	Logout(w, req)

	if w.Result().StatusCode != http.StatusFound {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusFound)
	}

	result := w.Result().Cookies()[0].Expires
	expected := time.Now().AddDate(0, 0, -1)
	if !result.Before(time.Now()) {
		t.Errorf("\nwrong Expiration date\ngot: %s\nexpected: %s",
			result.String(), expected.String())
	}
}