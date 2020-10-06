package jsonRealisation

import (
	"encoding/json"
	"io"
)

type UserJSON struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
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
