package signup

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/domain/entity/jsonRealisation"
	"server/handlers/authorization/utils"
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

	body := strings.NewReader(`{"email": "hound@psina.ru", "password1": "str", "password2": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusOK)
	}

	loc, _ := w.Result().Location()
	if loc.Path != utils.RootPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, utils.RootPage)
	}

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Errorf("no cookie")
	}

	if len(utils.UsersServerSession) == 0 {
		t.Errorf("no users created")
	}

	if len(utils.Sessions) == 0 {
		t.Errorf("no session")
	}
}

func TestSignupFailToComparePasswords(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := jsonRealisation.SignupJSON{
		Email:     "hound@psina.ru",
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
	if loc.Path != utils.SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, utils.SignupPage)
	}


}

func TestSignupFailWithNotFilledField(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := jsonRealisation.SignupJSON{
		Email:     "right",
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
	if loc.Path != utils.SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, utils.SignupPage)
	}
}

func TestSignupInvalidEmail(t *testing.T) {
	url := "http://example.com/api/"

	body := strings.NewReader(`{"email": "hound", "password1": "str", "password2": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	Signup(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("wrong StatusCode: got %d, expected %d",
			w.Code, http.StatusBadRequest)
	}

	cookies := w.Result().Cookies()
	if len(cookies) > 0 {
		t.Fatalf("cookie enabled")
	}

	loc, _ := w.Result().Location()
	if loc.Path != utils.SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, utils.SignupPage)
	}
}
