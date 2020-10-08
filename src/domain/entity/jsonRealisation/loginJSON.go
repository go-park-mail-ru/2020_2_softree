package jsonRealisation

import (
	"encoding/json"
	"io"
)

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

func (l *LoginJSON) GetEmail() string {
	return l.Email
}
