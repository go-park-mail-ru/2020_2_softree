package Entity

import (
	"encoding/json"
	"io"
)

type JSON interface {
	FillFields(io.ReadCloser) error
}

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

type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginJSON) FillFields(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&l)
	if err != nil {
		return err
	}

	return nil
}
