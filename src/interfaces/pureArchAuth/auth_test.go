package pureArchAuth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"server/src/domain/entity"
	"server/src/infrastructure/auth"
	"server/src/infrastructure/log"
	"server/src/infrastructure/persistence"
	"strings"
	"testing"
)

func TestAuthSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestAuthAuthenticateSuccess(req)

	testAuth.Auth(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.NotEmpty(t, w.Header().Get("Content-type"))
	assert.NotEmpty(t, w.Body)
}

func TestAuthFail(t *testing.T) {
	url := "http://127.0.0.1:8000/auth"
	body := strings.NewReader(`{"email": "yandex@mail.ru", "password": "str"}`)

	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()
	testAuth := createTestAuthAuthenticateFail()

	testAuth.Auth(w, req)
	assert.Empty(t, w.Header().Get("Content-type"))
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func createTestAuthAuthenticateSuccess(req *http.Request) *Authenticate {

}

func createTestAuthAuthenticateFail() *Authenticate {

}
