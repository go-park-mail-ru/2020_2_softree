package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/canal/pkg/domain/entity"
	"server/canal/pkg/infrastructure/mock"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/sessions"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginSuccess(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusOK, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestLogin_Fail_500(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginFail_500(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestLogin_Fail_400(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginFail_400(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestLogin_Fail_ErrorJSON(t *testing.T) {
	url := "http://127.0.0.1:8000/login"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createLoginFail_ErrorJSON(t, ctx)
	defer ctrl.Finish()

	testAuth.Login(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func createLoginSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Login(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{}, createExpectedUser(), http.Cookie{Name: name, Value: value}, nil)

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createLoginFail_500(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Login(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{Status: http.StatusInternalServerError}, entity.PublicUser{}, http.Cookie{}, errors.New("fail"))

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createLoginFail_400(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Login(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{Status: http.StatusBadRequest}, entity.PublicUser{}, http.Cookie{}, errors.New("fail"))

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createLoginFail_ErrorJSON(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Login(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{Status: http.StatusBadRequest,
			ErrorJSON: entity.ErrorJSON{NotEmpty: true, NonFieldError: []string{"Error"}}},
			entity.PublicUser{}, http.Cookie{},
			nil)

	return NewAuthentication(mockUser, mockAuth), ctrl
}
