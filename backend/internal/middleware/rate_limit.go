package middleware

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultRateLimitMaxKeys = 4096

// trustedProxyCount is how many trusted reverse-proxy hops sit in front of the
// app. It selects which entry of X-Forwarded-For is the real client IP for rate
// limiting: the value N hops from the right (0-indexed), i.e. the address the
// outermost trusted proxy observed.
//
// Because XFF is fully attacker-controlled, an untrusted value must never be
// used as a limiter key (a fresh IP per request bypasses the limit). The count
// is therefore explicit, not guessed from the header length:
//   - 0 (default, safe): ignore XFF/X-Real-IP entirely and key on RemoteAddr,
//     the immediate TCP peer. Correct when the app is exposed directly.
//   - N>=1: trust exactly N proxies and take the (N-1)-th value from the right.
//
// OPERATIONAL NOTE: set TRUSTED_PROXY_COUNT to the exact number of hops that
// prepend a *verified* client IP for this deployment (e.g. Cloud Run in front of
// the app is 1; a Vercel/CDN edge plus Cloud Run may be 2). Setting it larger
// than the real hop count re-opens the spoofing bypass; the default of 0 fails
// closed.
var trustedProxyCount = trustedProxyCountFromEnv()

func trustedProxyCountFromEnv() int {
	raw := strings.TrimSpace(os.Getenv("TRUSTED_PROXY_COUNT"))
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return 0
	}
	return n
}

type rateLimitEntry struct {
	windowStart time.Time
	count       int
}

type fixedWindowRateLimiter struct {
	limit   int
	window  time.Duration
	now     func() time.Time
	maxKeys int

	mu      sync.Mutex
	entries map[string]rateLimitEntry
}

func newFixedWindowRateLimiter(limit int, window time.Duration) *fixedWindowRateLimiter {
	return &fixedWindowRateLimiter{
		limit:   limit,
		window:  window,
		now:     time.Now,
		maxKeys: defaultRateLimitMaxKeys,
		entries: make(map[string]rateLimitEntry),
	}
}

func (l *fixedWindowRateLimiter) allow(key string) (bool, time.Duration) {
	if l == nil || l.limit <= 0 || l.window <= 0 {
		return true, 0
	}
	if key == "" {
		key = "unknown"
	}

	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.entries) >= l.maxKeys {
		l.pruneLocked(now)
	}

	entry, found := l.entries[key]
	if !found || !now.Before(entry.windowStart.Add(l.window)) {
		l.entries[key] = rateLimitEntry{
			windowStart: now,
			count:       1,
		}
		return true, 0
	}

	if entry.count >= l.limit {
		return false, entry.windowStart.Add(l.window).Sub(now)
	}

	entry.count++
	l.entries[key] = entry
	return true, 0
}

func (l *fixedWindowRateLimiter) pruneLocked(now time.Time) {
	for key, entry := range l.entries {
		if !now.Before(entry.windowStart.Add(l.window)) {
			delete(l.entries, key)
		}
	}
	for key := range l.entries {
		if len(l.entries) < l.maxKeys {
			return
		}
		delete(l.entries, key)
	}
}

// NewIPRateLimiter は、アプリから見えるクライアントIPごとにリクエスト数を制限する。
// Cloud Runへリクエストが到達した後の軽量なアプリ内ガードであり、
// 本格的なIP制限はCloud Armor/LBでCloud Runの手前に置く。
func NewIPRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newFixedWindowRateLimiter(limit, window)
	return newRateLimitMiddleware(limiter, clientIPFromRequest)
}

// NewAuthenticatedUserRateLimiter は、認証済みユーザーIDごとにリクエスト数を制限する。
// 認証ミドルウェアがリクエストコンテキストにユーザーIDを載せた後に配線する。
func NewAuthenticatedUserRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newFixedWindowRateLimiter(limit, window)
	return newRateLimitMiddleware(limiter, func(r *http.Request) string {
		userID := GetUserID(r.Context())
		if userID.IsZero() {
			return ""
		}
		return userID.String()
	})
}

func newRateLimitMiddleware(limiter *fixedWindowRateLimiter, keyFunc func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, retryAfter := limiter.allow(keyFunc(r))
			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(retryAfterSeconds(retryAfter)))
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func retryAfterSeconds(d time.Duration) int {
	if d <= 0 {
		return 1
	}
	seconds := int((d + time.Second - 1) / time.Second)
	if seconds < 1 {
		return 1
	}
	return seconds
}

func clientIPFromRequest(r *http.Request) string {
	return clientIPWithTrustedProxies(r, trustedProxyCount)
}

// clientIPWithTrustedProxies resolves the rate-limiting client IP given how many
// trusted reverse-proxy hops sit in front of the app. See trustedProxyCount for
// why forwarded headers are only trusted when trustedHops >= 1.
func clientIPWithTrustedProxies(r *http.Request, trustedHops int) string {
	if r == nil {
		return ""
	}

	// When no proxy is trusted, forwarded headers are attacker-controlled and
	// must be ignored; the immediate TCP peer is the only reliable key.
	if trustedHops <= 0 {
		return remoteAddrIP(r)
	}

	// XFF is "client, proxy1, proxy2, ...": the value trustedHops from the right
	// is what the outermost trusted proxy observed as the client. A shorter list
	// than expected means the request did not traverse the trusted chain, so we
	// fail closed to the direct peer rather than trusting a client-supplied value.
	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if idx := len(parts) - trustedHops; idx >= 0 {
			if ip := net.ParseIP(strings.TrimSpace(parts[idx])); ip != nil {
				return ip.String()
			}
		}
	}

	// X-Real-IP is a single value set by the nearest trusted proxy.
	if ip := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); ip != nil {
		return ip.String()
	}

	return remoteAddrIP(r)
}

func remoteAddrIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	if ip := net.ParseIP(strings.TrimSpace(r.RemoteAddr)); ip != nil {
		return ip.String()
	}
	return ""
}
