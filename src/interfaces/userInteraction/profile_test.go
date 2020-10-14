package userInteraction

import (
	"net/http"
	"net/http/httptest"
	"server/src/domain/entity"
	"server/src/infrastructure/security"
	"server/src/interfaces/authorization/utils"
	"strings"
	"testing"
)

func TestUpdateUserWithoutCookie(t *testing.T) {
	url := "http://127.0.0.1:8000/user"

	body := strings.NewReader(`{"email": "hound@psina.ru", "password1": "str", "password2": "str"}`)
	req := httptest.NewRequest("POST", url, body)
	w := httptest.NewRecorder()

	UpdateUserPartly(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("wrong status code\nexpected: %d\nreceived: %d", http.StatusUnauthorized, w.Code)
	}
}

var testEmail = "hound@psina.ru"
var testPassword = "123"
var cookie, _ = security.MakeCookie()

func makeTestData() {
	utils.UsersServerSession[testEmail], _ = security.MakeShieldedHash(testPassword)
	entity.Users = append(entity.Users, entity.PublicUser{Email: testEmail})
	utils.Sessions[testEmail] = cookie.Value
}

func TestUpdateUserSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/user"

	body := strings.NewReader(
		`{"oldPassword": "","newPassword1": "","newPassword2": "","avatar": "1234"}`,
	)
	req := httptest.NewRequest("PATCH", url, body)
	w := httptest.ResponseRecorder{}

	makeTestData()
	req.Header.Set("Cookie", cookie.String())
	UpdateUserPartly(&w, req)

	for _, obj := range entity.Users {
		if obj.Avatar == "" {
			t.Fatal("fail to change avatar")
		}
	}
}

func TestUpdateUserFailWrongOldPassword(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"

	body := strings.NewReader(
		`{"oldPassword": "12","newPassword1": "1234","newPassword2": "1234","avatar": ""}`,
	)
	req := httptest.NewRequest("PATCH", url, body)
	w := httptest.ResponseRecorder{}

	makeTestData()
	req.Header.Set("Cookie", cookie.String())
	UpdatePassword(&w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("wrong status code\nexpected: %d\nreceived: %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserFailWrongNewPassword(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"

	body := strings.NewReader(
		`{"oldPassword": "123","newPassword1": "1234","newPassword2": "123","avatar": ""}`,
	)
	req := httptest.NewRequest("PATCH", url, body)
	w := httptest.ResponseRecorder{}

	makeTestData()
	req.Header.Set("Cookie", cookie.String())
	UpdatePassword(&w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("wrong status code\nexpected: %d\nreceived: %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdatePasswordWithoutCookie(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"

	body := strings.NewReader(
		`{"oldPassword": "12","newPassword1": "1234","newPassword2": "1234","avatar": ""}`,
	)
	req := httptest.NewRequest("PATCH", url, body)
	w := httptest.NewRecorder()

	UpdatePassword(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("wrong status code\nexpected: %d\nreceived: %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdatePasswordSuccess(t *testing.T) {
	url := "http://127.0.0.1:8000/change-password"

	body := strings.NewReader(
		`{"oldPassword": "123","newPassword1": "1234","newPassword2": "1234","avatar": ""}`,
	)
	req := httptest.NewRequest("PATCH", url, body)
	w := httptest.ResponseRecorder{}

	makeTestData()
	req.Header.Set("Cookie", cookie.String())
	UpdatePassword(&w, req)

	passwordHash, _ := security.MakeShieldedHash("1234")
	if passwordHash != utils.UsersServerSession[testEmail] {
		t.Fatal("fail to update password")
	}

	if w.Code != http.StatusOK {
		t.Fatalf("wrong status code\nexpected: %d\nreceived: %d", http.StatusOK, w.Code)
	}
}
