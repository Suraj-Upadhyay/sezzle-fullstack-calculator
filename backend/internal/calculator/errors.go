package calculator

import (
	"errors"
	"fmt"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
)

var (
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrOperandCount         = errors.New("wrong number of operands")
	ErrInvalidOperand       = errors.New("operand is not a finite number")
	ErrDivisionByZero       = errors.New("division by zero")
	ErrUndefinedResult      = errors.New("result is not a real number")
	ErrResultOutOfRange     = errors.New("result is out of range")
	ErrInvalidExpression    = errors.New("invalid expression")
)

// SyntaxError reports where an expression stops making sense. Position
// counts characters from 1.
type SyntaxError struct {
	Position int
	Reason   string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("%s at position %d", e.Reason, e.Position)
}

func (e *SyntaxError) Unwrap() error {
	return ErrInvalidExpression
}

func (e *SyntaxError) APIDetail() map[string]any {
	return map[string]any{"position": e.Position, "reason": e.Reason}
}

var errorCodes = map[error]string{
	ErrUnsupportedOperation: apierror.CodeUnsupportedOperation,
	ErrOperandCount:         apierror.CodeOperandCount,
	ErrInvalidOperand:       apierror.CodeInvalidOperand,
	ErrDivisionByZero:       apierror.CodeDivisionByZero,
	ErrUndefinedResult:      apierror.CodeUndefinedResult,
	ErrResultOutOfRange:     apierror.CodeResultOutOfRange,
	ErrInvalidExpression:    apierror.CodeInvalidExpression,
}

func toAPIError(err error) error {
	for domainErr, code := range errorCodes {
		if errors.Is(err, domainErr) {
			return apierror.New(code, err)
		}
	}
	return err
}
