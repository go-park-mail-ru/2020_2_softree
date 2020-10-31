package profile

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
	testAuth, ctrl := createTestUpdateUserAuthenticateSuccess(t, entity.User{Avatar: "fake_image"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-type"))
	require.NotEmpty(t, w.Body)
}

func TestUpdateUserPasswordSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"
	body := strings.NewReader(`{"password": "fake_password"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createTestUpdateUserAuthenticateSuccess(t, entity.User{Password: "fake_password"})
	defer ctrl.Finish()

	createContext(&req)
	testAuth.UpdateUser(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Header().Get("Content-type"))
	require.NotEmpty(t, w.Body)
}

func createTestUpdateUserAuthenticateSuccess(t *testing.T, toUpdate entity.User) (*Profile, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	expectedUser := createExpectedUser()

	var id uint64 = 1
	mockUser := mocks.NewUserRepositoryForMock(ctrl)
	mockUser.EXPECT().UpdateUser(id, toUpdate).Return(expectedUser, nil)

	servicesDB := application.NewUserApp(mockUser)
	servicesAuth := auth.NewMemAuth()
	servicesCookie := auth.NewToken()
	servicesLog := log.NewLogrusLogger()

	return NewProfile(*servicesDB, servicesAuth, servicesCookie, servicesLog), ctrl
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

func createContext(req **http.Request) {
	ctx := context.WithValue((*req).Context(), "id", uint64(1))
	*req = (*req).Clone(ctx)
}
