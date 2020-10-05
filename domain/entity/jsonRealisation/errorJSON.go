package jsonRealisation

import (
	"encoding/json"
	"io"
)

type ErrorJSON struct {
	Name          []string `json:"name,omitempty"`
	Email         []string `json:"email,omitempty"`
	Password      []string `json:"password,omitempty"`
	NonFieldError []string `json:"non_field_error,omitempty"`
	NotEmpty      bool `json:"-"`
}

func (l *ErrorJSON) FillFields(body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&l)
	if err != nil {
		return err
	}

	return nil
}
