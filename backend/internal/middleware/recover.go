package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/apierror"
	"github.com/Suraj-Upadhyay/calculator/backend/internal/httputil"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}
			if err, ok := rec.(error); ok &&
				errors.Is(err, http.ErrAbortHandler) {
				panic(rec)
			}
			slog.Error("panic serving request", "panic", rec,
				"path", r.URL.Path, "stack", string(debug.Stack()))
			httputil.RenderError(w, apierror.New(apierror.CodeInternal,
				fmt.Errorf("panic: %v", rec)))
		}()
		next.ServeHTTP(w, r)
	})
}
