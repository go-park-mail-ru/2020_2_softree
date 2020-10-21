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
	servicesDB := persistence.NewUserRepository("db")
	servicesAuth := auth.NewMemAuth("auth")
	servicesCookie := auth.NewToken("token")

	testAuth := NewAuthenticate(servicesDB, servicesAuth, servicesCookie)

	url := "http://127.0.0.1:8000/signup"
	body := strings.NewReader(`{"email": "hound@psina.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	testAuth.Signup(w, req)

	assert.NotEmpty(t, persistence.Users)
	assert.NotEmpty(t, auth.Sessions)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
}
