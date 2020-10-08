package jsonRealisation

import (
	"encoding/json"
	"io"
)

type UserJSON struct {
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	OldPassword  string `json:"oldPassword"`
	NewPassword1 string `json:"newPassword1"`
	NewPassword2 string `json:"newPassword2"`
}

func (s *UserJSON) FillFields(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&s)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserJSON) GetEmail() string {
	return s.Email
}