package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterBlocksAfterLimit(t *testing.T) {
	limiter := newFixedWindowRateLimiter(2, time.Minute)
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }

	handler := newRateLimitMiddleware(limiter, func(*http.Request) string { return "client-1" })(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusNoContent {
			t.Fatalf("request %d status = %d, want %d", i+1, w.Code, http.StatusNoContent)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
	if got := w.Header().Get("Retry-After"); got != "60" {
		t.Fatalf("Retry-After = %q, want 60", got)
	}
}

func TestRateLimiterResetsAfterWindow(t *testing.T) {
	limiter := newFixedWindowRateLimiter(1, time.Minute)
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	limiter.now = func() time.Time { return now }

	if ok, _ := limiter.allow("client-1"); !ok {
		t.Fatal("first request should be allowed")
	}
	if ok, _ := limiter.allow("client-1"); ok {
		t.Fatal("second request in same window should be blocked")
	}

	now = now.Add(time.Minute)
	if ok, _ := limiter.allow("client-1"); !ok {
		t.Fatal("request after window should be allowed")
	}
}

func TestClientIPFromRequestUsesRightmostForwardedForIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "spoofed, 203.0.113.10, 2001:db8::1")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPFromRequest(req); got != "2001:db8::1" {
		t.Fatalf("clientIPFromRequest = %q, want rightmost valid X-Forwarded-For IP", got)
	}
}

func TestClientIPFromRequestFallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.20:12345"

	if got := clientIPFromRequest(req); got != "203.0.113.20" {
		t.Fatalf("clientIPFromRequest = %q, want RemoteAddr host", got)
	}
}
