package authorization

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	authMock "server/src/authorization/pkg/infrastructure/mock"
	session "server/src/authorization/pkg/session/gen"
	profileMock "server/src/profile/pkg/infrastructure/mock"
	profile "server/src/profile/pkg/profile/gen"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSignup_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupSuccess(t, &ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestSignup_FailUserExist(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupFailUserExist(t, &ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.NotEmpty(t, w.Header().Get("Content-Type"))
	require.NotEmpty(t, w.Body)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
}

func TestSignup_FailUserExistFailDB(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupFailUserExistFailDB(t, &ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestSignup_FailEmail(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFail(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestSignup_FailEmptyPassword(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth, ctrl := createSignupFail(t)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.NotEmpty(t, w.Body)
	require.NotEmpty(t, w.Header().Get("Content-type"))
}

func TestSignup_FailBcrypt(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": %s, \"password\": %s}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupFailBcrypt(t, &ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func createSignupSuccess(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email}).
		Return(&profile.Check{Existence: false}, nil)
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Times(1).
		Return(createExpectedUser(), nil)

	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)
	mockAuth.EXPECT().
		Create(ctx, &session.Session{Id: id, SessionId: value}).
		Return(&session.UserID{Id: id}, nil)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createSignupFail(t *testing.T) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := profileMock.NewProfileMock(ctrl)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createSignupFailBcrypt(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email}).
		Return(&profile.Check{Existence: false}, nil)
	mockUser.EXPECT().
		SaveUser(ctx, &profile.User{Email: email, Password: password}).
		Times(1).
		Return(nil, errors.New("createSignupFailBcrypt"))

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createSignupFailUserExist(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email}).
		Return(&profile.Check{Existence: true}, nil)

	return NewAuthenticate(mockUser, mockAuth), ctrl
}

func createSignupFailUserExistFailDB(t *testing.T, ctx *context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockAuth := authMock.NewAuthRepositoryForMock(ctrl)

	mockUser := profileMock.NewProfileMock(ctrl)
	mockUser.EXPECT().
		CheckExistence(ctx, &profile.User{Email: email}).
		Return(&profile.Check{Existence: false}, errors.New("createSignupFailUserExistFailDB"))

	return NewAuthenticate(mockUser, mockAuth), ctrl
}
