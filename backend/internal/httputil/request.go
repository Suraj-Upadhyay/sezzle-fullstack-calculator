package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
)

const MaxBodyBytes = 4 << 10

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// DecodeJSON reads exactly one JSON object of type T from the request body,
// rejecting unknown fields and trailing data, then validates it.
func DecodeJSON[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var payload T

	mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mediaType != "application/json" {
		return payload, apierror.New(apierror.CodeUnsupportedMediaType, nil)
	}

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, MaxBodyBytes))
	dec.DisallowUnknownFields()

	if err := dec.Decode(&payload); err != nil {
		return payload, decodeError(err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return payload, apierror.New(apierror.CodeInvalidJSON,
			errors.New("body must contain a single JSON object"))
	}

	if err := validate.Struct(payload); err != nil {
		return payload, validationError(err)
	}
	return payload, nil
}

func decodeError(err error) error {
	var maxBytes *http.MaxBytesError
	var typeErr *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		return apierror.New(apierror.CodeEmptyBody, err)
	case errors.As(err, &maxBytes):
		return apierror.New(apierror.CodeBodyTooLarge, err)
	case errors.As(err, &typeErr) && typeErr.Field != "":
		return apierror.New(apierror.CodeValidation, apierror.FieldErrors{
			typeErr.Field: "must be " + jsonType(typeErr.Type),
		})
	default:
		return apierror.New(apierror.CodeInvalidJSON, err)
	}
}

func jsonType(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int64:
		return "a number"
	case reflect.Slice, reflect.Array:
		return "an array"
	case reflect.Struct, reflect.Map:
		return "an object"
	default:
		return "a " + t.Kind().String()
	}
}

func validationError(err error) error {
	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		return err
	}

	fields := apierror.FieldErrors{}
	for _, f := range invalid {
		_, name, _ := strings.Cut(f.Namespace(), ".")
		fields[name] = fieldMessage(f)
	}
	return apierror.New(apierror.CodeValidation, fields)
}

func fieldMessage(f validator.FieldError) string {
	switch f.Tag() {
	case "required":
		return "is required"
	case "max":
		return "must be at most " + f.Param() + " characters"
	default:
		return "failed " + f.Tag() + " validation"
	}
}
