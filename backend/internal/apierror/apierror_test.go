package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewUsesRegisteredSpec(t *testing.T) {
	cause := errors.New("cause")
	err := New(CodeDivisionByZero, cause)

	if err.Status != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", err.Status,
			http.StatusUnprocessableEntity)
	}
	if !errors.Is(err, cause) {
		t.Error("error does not unwrap to its cause")
	}
}

func TestNewFallsBackToInternal(t *testing.T) {
	err := New("err_not_registered", nil)

	if err.Code != CodeInternal || err.Status != 500 {
		t.Errorf("got %s/%d, want %s/500", err.Code, err.Status,
			CodeInternal)
	}
}

func TestDetailIncludesFieldErrors(t *testing.T) {
	fields := FieldErrors{"operation": "is required"}
	err := New(CodeValidation, fields)

	got, ok := err.Detail()["fields"].(FieldErrors)
	if !ok || got["operation"] != "is required" {
		t.Errorf("detail = %v, want fields %v", err.Detail(), fields)
	}
}

func TestDetailIsEmptyWithoutFieldErrors(t *testing.T) {
	if d := New(CodeInternal, errors.New("x")).Detail(); len(d) != 0 {
		t.Errorf("detail = %v, want empty", d)
	}
}

func TestEveryCodeIsRegistered(t *testing.T) {
	codes := []string{
		CodeInvalidJSON, CodeEmptyBody, CodeBodyTooLarge,
		CodeUnsupportedMediaType, CodeValidation,
		CodeUnsupportedOperation, CodeOperandCount, CodeInvalidOperand,
		CodeDivisionByZero, CodeUndefinedResult, CodeResultOutOfRange,
		CodeInvalidExpression, CodeInternal,
	}
	for _, code := range codes {
		if _, ok := registry[code]; !ok {
			t.Errorf("%s has no registry entry", code)
		}
	}
}
