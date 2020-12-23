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

func TestSignup_Success(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupSuccess(t, ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.NotEmpty(t, w.Body)
	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
	require.NotEmpty(t, w.Result().Cookies())
}

func TestSignup_Fail(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupFail(t, ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.Empty(t, w.Header().Get("Content-Type"))
	require.Empty(t, w.Body)
	require.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	require.Empty(t, w.Result().Cookies())
}

func TestSignup_Fail_ErrorJSON(t *testing.T) {
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(fmt.Sprintf("{\"email\": \"%s\", \"password\": \"%s\"}", email, password))

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	ctx := req.Context()
	testAuth, ctrl := createSignupFail_ErrorJSON(t, ctx)
	defer ctrl.Finish()

	testAuth.Signup(w, req)

	require.NotEmpty(t, w.Body)
	require.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	require.Empty(t, w.Result().Cookies())
}

func createSignupSuccess(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Signup(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{}, createExpectedUser(), http.Cookie{Name: name, Value: value}, nil)

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createSignupFail(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Signup(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{Status: http.StatusInternalServerError}, createExpectedUser(), http.Cookie{}, errors.New("fail"))

	return NewAuthentication(mockUser, mockAuth), ctrl
}

func createSignupFail_ErrorJSON(t *testing.T, ctx context.Context) (*Authentication, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUser := mock.NewMockProfileLogic(ctrl)
	mockAuth := mock.NewMockAuthLogic(ctrl)

	mockAuth.EXPECT().
		Signup(ctx, entity.User{Email: email, Password: password}).
		Return(entity.Description{Status: http.StatusBadRequest,
			ErrorJSON: entity.ErrorJSON{NotEmpty: true, NonFieldError: []string{"Error"}}},
			entity.PublicUser{},
			http.Cookie{},
			nil)

	return NewAuthentication(mockUser, mockAuth), ctrl
}
