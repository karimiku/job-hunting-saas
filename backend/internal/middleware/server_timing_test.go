package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewServerTimingAddsAppMetric(t *testing.T) {
	handler := NewServerTiming()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/health", nil))

	serverTiming := w.Result().Header.Get("Server-Timing")
	if !strings.Contains(serverTiming, "app;dur=") {
		t.Fatalf("Server-Timing = %q, want app metric", serverTiming)
	}
}

func TestAddServerTimingMetricAddsNamedMetric(t *testing.T) {
	handler := NewServerTiming()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AddServerTimingMetric(r.Context(), "auth test", 1500*time.Microsecond)
		w.WriteHeader(http.StatusNoContent)
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil))

	serverTiming := w.Result().Header.Get("Server-Timing")
	if !strings.Contains(serverTiming, "auth_test;dur=1.5") {
		t.Fatalf("Server-Timing = %q, want sanitized auth_test metric", serverTiming)
	}
}

func TestNewServerTimingMergesExistingMetrics(t *testing.T) {
	handler := NewServerTiming()(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Server-Timing", "existing;dur=2.0")
		w.WriteHeader(http.StatusNoContent)
	}))

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/auth/session", nil))

	serverTiming := w.Result().Header.Get("Server-Timing")
	if !strings.Contains(serverTiming, "existing;dur=2.0") || !strings.Contains(serverTiming, "app;dur=") {
		t.Fatalf("Server-Timing = %q, want existing and app metrics", serverTiming)
	}
}
