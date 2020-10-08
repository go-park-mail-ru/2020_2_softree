package jsonRealisation

import (
	"encoding/json"
	"io"
)

type SignupJSON struct {
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

func (s *SignupJSON) FillFields(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&s)
	if err != nil {
		return err
	}

	return nil
}

func (s *SignupJSON) GetEmail() string {
	return s.Email
}
