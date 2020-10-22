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

func createTestAuthAuthenticateSuccess(req *http.Request) *Authenticate {
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

func createTestAuthAuthenticateFail(req *http.Request) *Authenticate {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")

	return NewAuthenticate(servicesDB, servicesAuth, servicesCookie)
}

func TestAuthSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestAuthAuthenticateSuccess(req)

	testAuth.Auth(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, auth.Sessions)
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, persistence.Users)
}

func TestAuthFail(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestAuthAuthenticateFail(req)

	testAuth.Auth(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

