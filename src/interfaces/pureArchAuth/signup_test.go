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

func TestAuthenticate_SignupSuccess(t *testing.T) {
	testAuth := createTestSignupAuthenticate(t)
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.Signup(w, req)

	assert.Empty(t, w.Header().Get("Content-type"))
	assert.Empty(t, w.Body)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func TestAuthenticate_SignupFailEmail(t *testing.T) {
	testAuth := createTestSignupAuthenticate(t)
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.Signup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestAuthenticate_SignupFailEmptyPassword(t *testing.T) {
	testAuth := createTestSignupAuthenticate(t)
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.Signup(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
}

func createTestSignupAuthenticate(t *testing.T) *Authenticate {
	userToSave := entity.User{
		Email: "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(userToSave.Password)
	expectedUser := entity.User{
		ID: 1,
		Email: userToSave.Email,
		Password: password,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUser := userMock.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().SaveUser(userToSave).Times(1).Return(expectedUser, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := auth.NewMemAuth()
	servicesCookie := auth.NewToken("token")
	servicesLog := log.NewLogrusLogger()

	return NewAuthenticate(*servicesDB, servicesAuth, servicesCookie, servicesLog)
}
