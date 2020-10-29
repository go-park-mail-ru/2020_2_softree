package profile

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
	"server/src/infrastructure/persistence"
	"strings"
	"testing"
)

func TestUpdateUserAvatarSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/user"
	body := strings.NewReader(`{"avatar": "fake_image"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateSuccess(req)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, auth.Sessions)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, persistence.Users)
}

func TestUpdateUserPasswordSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateSuccess(req)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, auth.Sessions)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, persistence.Users)
}

func TestUpdateUserFail(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "new_password"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateFail()

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Empty(t, w.Header().Get("Content-type"))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func createTestUpdateUserAuthenticateSuccess(req *http.Request) *Profile {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")
	servicesLog := log.NewLogrusLogger()

	cookie, _ := auth.CreateCookie()
	user := entity.User{Email: "yandex@mail.ru", Password: "str", Avatar: "some"}

	servicesDB.SaveUser(user)
	servicesAuth.CreateAuth(user.ID, cookie.Value)

	req.AddCookie(&cookie)
	return NewProfile(servicesDB, servicesAuth, servicesCookie, servicesLog)
}

func createTestUpdateUserAuthenticateFail() *Profile {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")
	servicesLog := log.NewLogrusLogger()

	return NewProfile(servicesDB, servicesAuth, servicesCookie, servicesLog)
}

