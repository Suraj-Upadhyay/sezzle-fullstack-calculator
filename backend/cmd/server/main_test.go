package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Suraj-Upadhyay/calculator/backend/internal/config"
)

func TestRegisterRoutesServesAPIAndStaticFiles(t *testing.T) {
	dir := t.TempDir()
	index := filepath.Join(dir, "index.html")
	if err := os.WriteFile(index, []byte("<h1>calc</h1>"), 0o600); err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux, config.Config{StaticDir: dir})

	tests := []struct {
		method, path, body string
		status             int
		contains           string
	}{
		{"GET", "/", "", http.StatusOK, "<h1>calc</h1>"},
		{"GET", "/healthz", "", http.StatusOK, "ok"},
		{"POST", "/api/v1/calculate", `{"operation":"add","operands":[1,2]}`,
			http.StatusOK, `"result":3`},
	}
	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path,
			strings.NewReader(tt.body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != tt.status ||
			!strings.Contains(rec.Body.String(), tt.contains) {
			t.Errorf("%s %s = %d %q", tt.method, tt.path, rec.Code, rec.Body)
		}
	}
}

func TestRegisterRoutesWithoutStaticDir(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux, config.Config{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}

func TestRunShutsDownWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- run(ctx, config.Config{Port: "0", ShutdownTimeout: time.Second})
	}()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("run returned %v, want nil", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("run did not return after cancellation")
	}
}

func TestRunReportsListenErrors(t *testing.T) {
	err := run(context.Background(), config.Config{Port: "-1"})
	if err == nil {
		t.Error("run returned nil for an invalid port")
	}
}
