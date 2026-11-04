package middleware

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func captureLogs(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func serve(h http.Handler) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
	return rec
}

func TestRecoverRendersInternalError(t *testing.T) {
	logs := captureLogs(t)
	h := Wrap(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}))

	rec := serve(h)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"error":"err_internal"`) {
		t.Errorf("body = %s", rec.Body)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Error("panic value leaked to the client")
	}
	if !strings.Contains(logs.String(), "status=500") {
		t.Errorf("logger did not record the 500: %s", logs)
	}
}

func TestRecoverRepanicsOnAbortHandler(t *testing.T) {
	h := Recover(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	defer func() {
		err, _ := recover().(error)
		if !errors.Is(err, http.ErrAbortHandler) {
			t.Error("ErrAbortHandler was swallowed")
		}
	}()
	serve(h)
}

func TestLoggerRecordsStatus(t *testing.T) {
	logs := captureLogs(t)
	h := Logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	if rec := serve(h); rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want 418", rec.Code)
	}
	for _, want := range []string{"method=GET", "path=/x", "status=418"} {
		if !strings.Contains(logs.String(), want) {
			t.Errorf("log %q missing %q", logs, want)
		}
	}
}

func TestLoggerDefaultsToOK(t *testing.T) {
	logs := captureLogs(t)
	h := Logger(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("hi"))
	}))

	serve(h)

	if !strings.Contains(logs.String(), "status=200") {
		t.Errorf("log %q missing status=200", logs)
	}
}
