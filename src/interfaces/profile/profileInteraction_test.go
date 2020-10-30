package profile

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/application"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
	mocks "server/src/infrastructure/mock"
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

	cookie := http.Cookie{Value: "value"}
	req.AddCookie(&cookie)

	update := testAuth.Auth(testAuth.UpdateUser)
	update(w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)

	fmt.Println(w.Body)
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedUser := createExpectedUser()
	updatedUserFields := entity.User{Avatar: "fake_image"}

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUser(id, updatedUserFields).Times(1).Return(expectedUser, nil)

	mockAuth := mocks.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().CheckAuth("value").Times(1).Return(id, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := auth.NewMemAuth()
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, servicesAuth, servicesCookie, servicesLog)
}

func createExpectedUser() (expected entity.User) {
	toSave := entity.User{
		Email: "hound@psina.ru",
		Password: "str",
	}
	password, _ := security.MakeShieldedPassword(toSave.Password)
	expected = entity.User{
		ID: 1,
		Email: toSave.Email,
		Password: password,
		Avatar: "fake_image",
	}

	return
}
