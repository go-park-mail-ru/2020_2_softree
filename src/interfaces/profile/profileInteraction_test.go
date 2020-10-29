package profile

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

func TestUpdateUserAvatarSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/user"
	body := strings.NewReader(`{"avatar": "fake_image"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateSuccess(t)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)
}

func TestUpdateUserPasswordSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateSuccess(t)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)
}

func TestUpdateUserFail(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "new_password"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestUpdateUserAuthenticateSuccess(t)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Empty(t, w.Header().Get("Content-type"))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func createTestUpdateUserAuthenticateSuccess(t *testing.T) *Profile {
	userToSave := entity.User{
		Email: "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedHash(userToSave.Password)
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

	return NewProfile(*servicesDB, servicesAuth, servicesCookie, servicesLog)
}
