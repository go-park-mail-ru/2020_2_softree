package pureArchAuth

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
	userMock "server/src/infrastructure/mock"
	"server/src/infrastructure/security"
	"strings"
	"testing"
)

func TestLogoutSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestLogoutAuthenticateSuccess(t)

	testAuth.Logout(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode)
}

func TestLogoutFail(t *testing.T) {
	url := "http://127.0.0.1:8000/logout"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestLogoutAuthenticateSuccess(t)

	testAuth.Logout(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode)
}

func createTestLogoutAuthenticateSuccess(t *testing.T) *Authenticate {
	userToSave := entity.User{
		Email:    "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(userToSave.Password)
	expectedUser := entity.User{
		ID:       1,
		Email:    userToSave.Email,
		Password: password,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := userMock.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	memAuth := auth.NewMemAuth()

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := application.NewUserAuth(memAuth)
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, *servicesAuth, servicesCookie, servicesLog)
}
