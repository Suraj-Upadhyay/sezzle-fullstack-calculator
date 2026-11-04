package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
)

type payload struct {
	Name  string   `json:"name" validate:"required"`
	Tags  []string `json:"tags"`
	Inner struct {
		N int `json:"n"`
	} `json:"inner"`
}

func decode(body string) error {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	_, err := DecodeJSON[payload](httptest.NewRecorder(), req)
	return err
}

func TestDecodeJSONReportsWireTypes(t *testing.T) {
	tests := map[string]string{
		`{"name":"a","tags":"x"}`:       "must be an array",
		`{"name":"a","inner":1}`:        "must be an object",
		`{"name":"a","tags":[1]}`:       "must be a string",
		`{"name":"a","inner":{"n":""}}`: "must be a number",
	}
	for body, want := range tests {
		var apiErr *apierror.Error
		if err := decode(body); !errors.As(err, &apiErr) {
			t.Fatalf("%s: err = %v", body, err)
		}
		fields, _ := apiErr.Detail()["fields"].(apierror.FieldErrors)
		for _, got := range fields {
			if got != want {
				t.Errorf("%s: got %q, want %q", body, got, want)
			}
		}
	}
}

func TestDecodeJSONAcceptsValidBody(t *testing.T) {
	if err := decode(`{"name":"a","tags":["x"]}`); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRegisterAppliesLastMiddlewareOutermost(t *testing.T) {
	var order []string
	tag := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter,
				r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}

	mux := http.NewServeMux()
	Register(mux, []Route{{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, _ *http.Request) error {
			return WriteJSON(w, http.StatusOK, map[string]bool{"ok": true})
		},
		Middleware: []func(http.Handler) http.Handler{
			tag("inner"), tag("outer"),
		},
	}})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if strings.Join(order, ",") != "outer,inner" {
		t.Errorf("order = %v, want [outer inner]", order)
	}
	var body map[string]bool
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil ||
		!body["ok"] {
		t.Errorf("body = %v, err = %v", body, err)
	}
}
