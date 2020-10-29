package pureArchAuth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/persistence"
	"strings"
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestLogoutAuthenticateSuccess(req)

	testAuth.Logout(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode)
	assert.Empty(t, auth.Sessions)
}

func TestLogoutFail(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestLogoutAuthenticateFail(req)

	testAuth.Logout(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode)
}

func createTestLogoutAuthenticateSuccess(req *http.Request) *Authenticate {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")

	cookie, _ := auth.CreateCookie()
	user := entity.User{Email: "yandex@mail.ru", Password: "str"}

	servicesDB.SaveUser(user)
	servicesAuth.CreateAuth(user.ID, cookie.Value)

	req.AddCookie(&cookie)
	return NewAuthenticate(servicesDB, servicesAuth, servicesCookie)
}

func createTestLogoutAuthenticateFail(req *http.Request) *Authenticate {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")

	return NewAuthenticate(servicesDB, servicesAuth, servicesCookie)
}
