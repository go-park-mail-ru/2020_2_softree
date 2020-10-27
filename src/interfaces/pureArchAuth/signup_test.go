package pureArchAuth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/persistence"
	"strings"
	"testing"
)

func TestAuthenticate_SignupSuccess(t *testing.T) {
	testAuth := createTestSignupAuthenticate()
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.PrepareUser(testAuth.SaveUser(testAuth.Signup))

	assert.NotEmpty(t, persistence.Users)
	assert.NotEmpty(t, auth.Sessions)
	assert.Empty(t, w.Header().Get("Content-type"))
	assert.Empty(t, w.Body)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}

func TestAuthenticate_SignupFailEmail(t *testing.T) {
	testAuth := createTestSignupAuthenticate()
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.PrepareUser(testAuth.SaveUser(testAuth.Signup))

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.Empty(t, persistence.Users)
	assert.Empty(t, auth.Sessions)
}

func TestAuthenticate_SignupFailEmptyPassword(t *testing.T) {
	testAuth := createTestSignupAuthenticate()
	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.PrepareUser(testAuth.SaveUser(testAuth.Signup))

	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	assert.NotEmpty(t, w.Body)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.Empty(t, persistence.Users)
	assert.Empty(t, auth.Sessions)
}

func createTestSignupAuthenticate() *Authenticate {
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")

	return NewAuthenticate(servicesDB, servicesAuth, servicesCookie)
}
