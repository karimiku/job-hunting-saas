package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultRateLimitMaxKeys = 4096

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

// NewIPRateLimiter limits requests by the best client IP visible to the app.
// It is a lightweight in-process guard; Cloud Armor/LB is still the stronger
// place to enforce IP policy before Cloud Run starts work.
func NewIPRateLimiter(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newFixedWindowRateLimiter(limit, window)
	return newRateLimitMiddleware(limiter, clientIPFromRequest)
}

// NewAuthenticatedUserRateLimiter limits requests by authenticated user ID.
// It should be wired after auth middleware has populated the request context.
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
	if r == nil {
		return ""
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		for i := len(parts) - 1; i >= 0; i-- {
			if ip := net.ParseIP(strings.TrimSpace(parts[i])); ip != nil {
				return ip.String()
			}
		}
	}

	if ip := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); ip != nil {
		return ip.String()
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	if ip := net.ParseIP(strings.TrimSpace(r.RemoteAddr)); ip != nil {
		return ip.String()
	}
	return ""
}
