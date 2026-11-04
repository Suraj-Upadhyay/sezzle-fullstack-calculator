package httputil

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
)

// APIHandler lets a handler return its error; the error is logged and
// rendered as JSON, and anything that is not an *apierror.Error is
// reported to the client as an internal error.
type APIHandler func(w http.ResponseWriter, r *http.Request) error

func (h APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err == nil {
		return
	}

	var apiErr *apierror.Error
	if !errors.As(err, &apiErr) {
		apiErr = apierror.New(apierror.CodeInternal, err)
	}

	level := slog.LevelWarn
	if apiErr.Status >= http.StatusInternalServerError {
		level = slog.LevelError
	}
	slog.Log(r.Context(), level, "request failed",
		"err", err, "method", r.Method, "path", r.URL.Path)

	RenderError(w, apiErr)
}

func RenderError(w http.ResponseWriter, e *apierror.Error) {
	body := map[string]any{
		"error":   e.Code,
		"message": e.Message,
	}
	for k, v := range e.Detail() {
		body[k] = v
	}
	if err := WriteJSON(w, e.Status, body); err != nil {
		slog.Error("writing error response", "err", err)
	}
}
