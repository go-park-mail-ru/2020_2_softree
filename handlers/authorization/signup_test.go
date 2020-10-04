package authorization

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/domain/Entity"
	"strings"
	"testing"
)

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
	if loc.Path != LoginPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, LoginPage)
	}
}

func TestSignupFailToComparePasswords(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := Entity.SignupJSON{
		Email: "right",
		Password1: "str",
		Password2: "ste",
	}
	body, _ := json.Marshal(jsonForBody)
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusBadRequest)
	}

	loc, _ := w.Result().Location()
	if loc.Path != SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, SignupPage)
	}
}

func TestSignupFailWithNotFilledField(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := Entity.SignupJSON{
		Email: "right",
		Password1: "str",
		Password2: "ste",
	}
	body, _ := json.Marshal(jsonForBody)
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("\nwrong StatusCode\ngot: %d\nexpected: %d",
			w.Code, http.StatusBadRequest)
	}

	loc, _ := w.Result().Location()
	if loc.Path != SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, SignupPage)
	}
}
