package Handlers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestCase struct {
	Email      string
	Password   string
	Response   string
	StatusCode int
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	expectedUser := make(map[string]string)

	h := md5.New()
	expectedEmail := "right"
	expectedPassword := "strongPassword"
	h.Write([]byte(expectedPassword))
	expectedUser[expectedEmail] = hex.EncodeToString(h.Sum(nil))
	if expectedUser[email] == hex.EncodeToString(h.Sum(nil)) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "resp": {"email": "right", "password": "123"}}`)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status": 400, "resp": {"email": "wrong", "password": "123"}}`)
	}
}

func TestSignupFailWithGET(t *testing.T) {
	url := "http://example.com/api/"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusBadRequest)
	}
}

func TestSignupSuccess(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "right", "password1": "str", "password2": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	loc, _ := w.Result().Location()
	if loc.Path != login {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, login)
	}
}

func TestSignupFailToComparePasswords(t *testing.T) {
	url := "http://example.com/api/"

	mapForJSON := map[string]string{
		"email":     "right",
		"password1": "str",
		"password2": "ste",
	}
	body, _ := json.Marshal(mapForJSON)
	req := httptest.NewRequest("POST", url, strings.NewReader(string(body)))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusBadRequest)
	}
}
