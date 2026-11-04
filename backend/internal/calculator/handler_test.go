package calculator

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

func newTestMux() *http.ServeMux {
	mux := http.NewServeMux()
	Register(mux)
	return mux
}

func post(
	t *testing.T, h http.Handler, contentType, body string,
) *httptest.ResponseRecorder {
	t.Helper()
	return postTo(t, h, "/api/v1/calculate", contentType, body)
}

func postTo(
	t *testing.T, h http.Handler, path, contentType, body string,
) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path,
		strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCalculateEndpoint(t *testing.T) {
	rec := post(t, newTestMux(), "application/json",
		`{"operation":"divide","operands":[10,4]}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content type = %q", ct)
	}

	var got CalculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	want := CalculateResponse{Divide, []float64{10, 4}, 2.5}
	if got.Operation != want.Operation || got.Result != want.Result ||
		len(got.Operands) != 2 {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestCalculateEndpointErrors(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		status      int
		code        string
		fields      map[string]string
	}{
		{
			name:        "wrong content type",
			contentType: "text/plain",
			body:        `{"operation":"add","operands":[1,2]}`,
			status:      http.StatusUnsupportedMediaType,
			code:        apierror.CodeUnsupportedMediaType,
		},
		{
			name:   "empty body",
			body:   ``,
			status: http.StatusBadRequest,
			code:   apierror.CodeEmptyBody,
		},
		{
			name:   "malformed json",
			body:   `{"operation":"add",`,
			status: http.StatusBadRequest,
			code:   apierror.CodeInvalidJSON,
		},
		{
			name:   "trailing data",
			body:   `{"operation":"add","operands":[1,2]} {}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeInvalidJSON,
		},
		{
			name:   "unknown field",
			body:   `{"operation":"add","operands":[1,2],"x":1}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeInvalidJSON,
		},
		{
			name:   "too large",
			body:   `{"operation":"` + strings.Repeat("a", 5000) + `"}`,
			status: http.StatusRequestEntityTooLarge,
			code:   apierror.CodeBodyTooLarge,
		},
		{
			name:   "number too large for float64",
			body:   `{"operation":"add","operands":[1e400,1]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeValidation,
			fields: map[string]string{"operands": "must be a number"},
		},
		{
			name:   "string operand",
			body:   `{"operation":"add","operands":["1",2]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeValidation,
			fields: map[string]string{"operands": "must be a number"},
		},
		{
			name:   "missing operation",
			body:   `{"operands":[1,2]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeValidation,
			fields: map[string]string{"operation": "is required"},
		},
		{
			name:   "missing operands",
			body:   `{"operation":"add"}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeValidation,
			fields: map[string]string{"operands": "is required"},
		},
		{
			name:   "null operand",
			body:   `{"operation":"add","operands":[1,null]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeValidation,
			fields: map[string]string{"operands[1]": "is required"},
		},
		{
			name:   "unsupported operation",
			body:   `{"operation":"modulo","operands":[1,2]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeUnsupportedOperation,
		},
		{
			name:   "wrong operand count",
			body:   `{"operation":"sqrt","operands":[1,2]}`,
			status: http.StatusBadRequest,
			code:   apierror.CodeOperandCount,
		},
		{
			name:   "division by zero",
			body:   `{"operation":"divide","operands":[1,0]}`,
			status: http.StatusUnprocessableEntity,
			code:   apierror.CodeDivisionByZero,
		},
		{
			name:   "square root of negative",
			body:   `{"operation":"sqrt","operands":[-4]}`,
			status: http.StatusUnprocessableEntity,
			code:   apierror.CodeUndefinedResult,
		},
		{
			name:   "overflow",
			body:   `{"operation":"power","operands":[10,400]}`,
			status: http.StatusUnprocessableEntity,
			code:   apierror.CodeResultOutOfRange,
		},
	}

	mux := newTestMux()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentType := tt.contentType
			if contentType == "" {
				contentType = "application/json; charset=utf-8"
			}
			rec := post(t, mux, contentType, tt.body)
			assertError(t, rec, tt.status, tt.code, tt.fields)
		})
	}
}

type failingCalculator struct{}

func (failingCalculator) Calculate(Calculation) (float64, error) {
	return 0, errors.New("unexpected failure")
}

func (failingCalculator) Evaluate(string) (float64, error) {
	return 0, errors.New("unexpected failure")
}

func TestUnexpectedErrorIsInternal(t *testing.T) {
	mux := http.NewServeMux()
	httputil.Register(mux, routes(NewHandler(failingCalculator{})))

	rec := post(t, mux, "application/json",
		`{"operation":"add","operands":[1,2]}`)

	assertError(t, rec, http.StatusInternalServerError,
		apierror.CodeInternal, nil)
	if strings.Contains(rec.Body.String(), "unexpected failure") {
		t.Error("internal error detail leaked to the client")
	}
}

func assertError(t *testing.T, rec *httptest.ResponseRecorder, status int,
	code string, fields map[string]string) {
	t.Helper()

	var body struct {
		Error   string            `json:"error"`
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decoding error body: %v", err)
	}
	if rec.Code != status || body.Error != code {
		t.Errorf("got %d %s, want %d %s", rec.Code, body.Error, status, code)
	}
	if body.Message == "" {
		t.Error("error has no message")
	}
	for field, msg := range fields {
		if body.Fields[field] != msg {
			t.Errorf("fields[%q] = %q, want %q (all: %v)",
				field, body.Fields[field], msg, body.Fields)
		}
	}
}

func TestEvaluateEndpoint(t *testing.T) {
	rec := postTo(t, newTestMux(), "/api/v1/evaluate", "application/json",
		`{"expression":"(2 + 3) × 4"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body %s", rec.Code, rec.Body)
	}
	var got EvaluateResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Expression != "(2 + 3) × 4" || got.Result != 20 {
		t.Errorf("got %+v", got)
	}
}

func TestEvaluateEndpointErrors(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
		code   string
	}{
		{"missing expression", `{}`, http.StatusBadRequest,
			apierror.CodeValidation},
		{"too long", `{"expression":"` + strings.Repeat("1", 501) + `"}`,
			http.StatusBadRequest, apierror.CodeValidation},
		{"syntax error", `{"expression":"2 +"}`, http.StatusBadRequest,
			apierror.CodeInvalidExpression},
		{"division by zero", `{"expression":"1/(2-2)"}`,
			http.StatusUnprocessableEntity, apierror.CodeDivisionByZero},
	}

	mux := newTestMux()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := postTo(t, mux, "/api/v1/evaluate", "application/json",
				tt.body)
			assertError(t, rec, tt.status, tt.code, nil)
		})
	}
}

func TestSyntaxErrorReportsPosition(t *testing.T) {
	rec := postTo(t, newTestMux(), "/api/v1/evaluate", "application/json",
		`{"expression":"(1+2"}`)

	var body struct {
		Position int    `json:"position"`
		Reason   string `json:"reason"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Position != 5 || body.Reason != `missing ")"` {
		t.Errorf("got %+v, want position 5 and missing )", body)
	}
}

func TestEvaluateUnexpectedErrorIsInternal(t *testing.T) {
	mux := http.NewServeMux()
	httputil.Register(mux, routes(NewHandler(failingCalculator{})))

	rec := postTo(t, mux, "/api/v1/evaluate", "application/json",
		`{"expression":"1+2"}`)

	assertError(t, rec, http.StatusInternalServerError,
		apierror.CodeInternal, nil)
}
