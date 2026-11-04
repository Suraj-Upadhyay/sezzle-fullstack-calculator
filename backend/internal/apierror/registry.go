package apierror

import "net/http"

type Spec struct {
	Status  int
	Message string
}

// Request errors
const (
	CodeInvalidJSON          = "err_invalid_json"
	CodeEmptyBody            = "err_empty_request_body"
	CodeBodyTooLarge         = "err_request_body_too_large"
	CodeUnsupportedMediaType = "err_unsupported_media_type"
	CodeValidation           = "err_validation"
)

// Calculation errors
const (
	CodeUnsupportedOperation = "err_unsupported_operation"
	CodeOperandCount         = "err_operand_count"
	CodeInvalidOperand       = "err_invalid_operand"
	CodeDivisionByZero       = "err_division_by_zero"
	CodeUndefinedResult      = "err_undefined_result"
	CodeResultOutOfRange     = "err_result_out_of_range"
	CodeInvalidExpression    = "err_invalid_expression"
)

// Server errors
const (
	CodeInternal = "err_internal"
)

var registry = map[string]Spec{
	CodeInvalidJSON: {
		http.StatusBadRequest, "request body is not valid JSON",
	},
	CodeEmptyBody: {
		http.StatusBadRequest, "request body is empty",
	},
	CodeBodyTooLarge: {
		http.StatusRequestEntityTooLarge, "request body is too large",
	},
	CodeUnsupportedMediaType: {
		http.StatusUnsupportedMediaType,
		"content type must be application/json",
	},
	CodeValidation: {
		http.StatusBadRequest, "request failed validation",
	},

	CodeUnsupportedOperation: {
		http.StatusBadRequest, "operation is not supported",
	},
	CodeOperandCount: {
		http.StatusBadRequest,
		"wrong number of operands for this operation",
	},
	CodeInvalidOperand: {
		http.StatusBadRequest, "operands must be finite numbers",
	},
	CodeDivisionByZero: {
		http.StatusUnprocessableEntity, "division by zero is undefined",
	},
	CodeUndefinedResult: {
		http.StatusUnprocessableEntity,
		"result is not a real number",
	},
	CodeResultOutOfRange: {
		http.StatusUnprocessableEntity,
		"result is too large to represent",
	},

	CodeInvalidExpression: {
		http.StatusBadRequest, "expression is not valid",
	},

	CodeInternal: {
		http.StatusInternalServerError, "internal server error",
	},
}
