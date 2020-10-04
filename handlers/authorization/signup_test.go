package authorization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"server/domain/entity/jsonRealisation"
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
	if loc.Path != RootPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, RootPage)
	}

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Errorf("no cookie ")
	}

	fmt.Println(len(UsersServerSession))
}

func TestSignupFailToComparePasswords(t *testing.T) {
	url := "http://example.com/api/"

	jsonForBody := jsonRealisation.SignupJSON{
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

	jsonForBody := jsonRealisation.SignupJSON{
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
	if loc.Path != SignupPage {
		t.Errorf("wrong Location: got %s, expected %s",
			loc.Path, SignupPage)
	}
}
