package entity

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

type ErrorJSON struct {
	Name []string `json:"name,omitempty"`
	Email    []string `json:"email,omitempty"`
	Password []string `json:"password,omitempty"`
}

func (l *ErrorJSON) FillFields(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&l)
	if err != nil {
		return err
	}

	return nil
}
