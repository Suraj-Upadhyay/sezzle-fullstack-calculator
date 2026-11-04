package apierror

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    string
	Status  int
	Message string
	Err     error
}

// New builds the error registered under code, falling back to CodeInternal
// for an unregistered code so internal detail never reaches a client.
func New(code string, err error) *Error {
	spec, ok := registry[code]
	if !ok {
		code, spec = CodeInternal, registry[CodeInternal]
	}
	return &Error{
		Code:    code,
		Status:  spec.Status,
		Message: spec.Message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %d: %s: %v", e.Code, e.Status, e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Detailer is implemented by errors that add fields to the error body.
type Detailer interface {
	APIDetail() map[string]any
}

func (e *Error) Detail() map[string]any {
	var d Detailer
	if errors.As(e.Err, &d) {
		return d.APIDetail()
	}
	return map[string]any{}
}

// FieldErrors maps a request field, named as it appears on the wire, to
// what is wrong with it.
type FieldErrors map[string]string

func (f FieldErrors) Error() string {
	return fmt.Sprintf("invalid fields: %v", map[string]string(f))
}

func (f FieldErrors) APIDetail() map[string]any {
	return map[string]any{"fields": f}
}
