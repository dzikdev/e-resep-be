package model

import (
	"fmt"
)

type ErrorKind string

const (
	Validation ErrorKind = "Validation Error"
	Type       ErrorKind = "Type Error"
	NotFound   ErrorKind = "Not Found"
	Unknown    ErrorKind = "Unknown Error"
)

// NewError return wrapped dynamic errors
func NewError(kind ErrorKind, msg string) error {
	return fmt.Errorf("%s: %s", string(kind), msg)
}
