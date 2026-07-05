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

func TestClientIPWithZeroTrustedProxiesIgnoresForwardedFor(t *testing.T) {
	// Default (0 trusted proxies): X-Forwarded-For is attacker-controlled and
	// must be ignored so a spoofed header cannot mint a fresh limiter key.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
	req.Header.Set("X-Real-IP", "9.9.9.9")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPWithTrustedProxies(req, 0); got != "10.0.0.1" {
		t.Fatalf("clientIPWithTrustedProxies(0) = %q, want RemoteAddr host 10.0.0.1", got)
	}
}

func TestClientIPWithOneTrustedProxyUsesForwardedForClient(t *testing.T) {
	// One trusted proxy prepends the real client. XFF = "client": the value one
	// hop from the right is the client the proxy observed.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPWithTrustedProxies(req, 1); got != "203.0.113.10" {
		t.Fatalf("clientIPWithTrustedProxies(1) = %q, want 203.0.113.10", got)
	}
}

func TestClientIPWithOneTrustedProxyIgnoresSpoofedLeftEntries(t *testing.T) {
	// An attacker prepending fake entries to XFF cannot shift the trusted index:
	// with 1 trusted proxy we always read the rightmost entry (what our proxy set).
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "1.1.1.1, 2.2.2.2, 203.0.113.99")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPWithTrustedProxies(req, 1); got != "203.0.113.99" {
		t.Fatalf("clientIPWithTrustedProxies(1) = %q, want rightmost 203.0.113.99", got)
	}
}

func TestClientIPWithTwoTrustedProxiesSelectsSecondFromRight(t *testing.T) {
	// Two trusted proxies: XFF = "client, proxy1"; index len-2 is the client.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.2")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPWithTrustedProxies(req, 2); got != "198.51.100.7" {
		t.Fatalf("clientIPWithTrustedProxies(2) = %q, want 198.51.100.7", got)
	}
}

func TestClientIPFailsClosedWhenForwardedForShorterThanTrustedHops(t *testing.T) {
	// Fewer XFF entries than trusted hops means the request did not traverse the
	// expected chain; fail closed to the direct peer rather than trust the value.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10")
	req.RemoteAddr = "10.0.0.1:12345"

	if got := clientIPWithTrustedProxies(req, 2); got != "10.0.0.1" {
		t.Fatalf("clientIPWithTrustedProxies(2) with short XFF = %q, want RemoteAddr 10.0.0.1", got)
	}
}

func TestClientIPFallsBackToRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "203.0.113.20:12345"

	if got := clientIPWithTrustedProxies(req, 1); got != "203.0.113.20" {
		t.Fatalf("clientIPWithTrustedProxies = %q, want RemoteAddr host", got)
	}
}

func TestTrustedProxyCountFromEnvDefaultsToZero(t *testing.T) {
	t.Setenv("TRUSTED_PROXY_COUNT", "")
	if got := trustedProxyCountFromEnv(); got != 0 {
		t.Fatalf("trustedProxyCountFromEnv() unset = %d, want 0", got)
	}
	t.Setenv("TRUSTED_PROXY_COUNT", "-3")
	if got := trustedProxyCountFromEnv(); got != 0 {
		t.Fatalf("trustedProxyCountFromEnv() negative = %d, want 0", got)
	}
	t.Setenv("TRUSTED_PROXY_COUNT", "not-a-number")
	if got := trustedProxyCountFromEnv(); got != 0 {
		t.Fatalf("trustedProxyCountFromEnv() invalid = %d, want 0", got)
	}
	t.Setenv("TRUSTED_PROXY_COUNT", "2")
	if got := trustedProxyCountFromEnv(); got != 2 {
		t.Fatalf("trustedProxyCountFromEnv() = %d, want 2", got)
	}
}
