package entity

import (
	"encoding/json"
	"io"
)

type ErrorJSON struct {
	Name          []string `json:"name,omitempty"`
	Email         []string `json:"email,omitempty"`
	Password      []string `json:"password,omitempty"`
	OldPassword   []string `json:"oldPassword,omitempty"`
	NonFieldError []string `json:"non_field_errors,omitempty"`
	NotEmpty      bool     `json:"-"`
}

func (l ErrorJSON) FillFields(body io.ReadCloser) error {
	err := json.NewDecoder(body).Decode(&l)
	if err != nil {
		return err
	}

	return nil
}
