package jsonRealisation

import "io"

type JSON interface {
	FillFields(io.ReadCloser) error
	GetEmail() string
}
