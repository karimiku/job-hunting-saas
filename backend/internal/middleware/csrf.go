package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// NewOriginGuard rejects unsafe browser requests unless Origin or Referer
// resolves to one of the explicitly allowed origins.
func NewOriginGuard(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := originSet(allowedOrigins)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isUnsafeMethod(r.Method) || requestOriginAllowed(r, allowed) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

// NewSessionCSRFProtection applies Origin/Referer checks only to unsafe
// requests authenticated by a browser session cookie. Bearer token clients are
// not subject to browser CSRF, so they are intentionally skipped.
func NewSessionCSRFProtection(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := originSet(allowedOrigins)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isUnsafeMethod(r.Method) ||
				GetAuthMethod(r.Context()) != AuthMethodSession ||
				requestOriginAllowed(r, allowed) {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func isUnsafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	default:
		return true
	}
}

func originSet(origins []string) map[string]struct{} {
	allowed := make(map[string]struct{}, len(origins))
	for _, raw := range origins {
		origin, ok := normalizeOrigin(raw)
		if ok {
			allowed[origin] = struct{}{}
		}
	}
	return allowed
}

func requestOriginAllowed(r *http.Request, allowed map[string]struct{}) bool {
	if len(allowed) == 0 {
		return false
	}
	if origin, ok := normalizeOrigin(r.Header.Get("Origin")); ok {
		_, found := allowed[origin]
		return found
	}
	if origin, ok := refererOrigin(r.Header.Get("Referer")); ok {
		_, found := allowed[origin]
		return found
	}
	return false
}

func refererOrigin(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	return originFromURL(parsed)
}

func normalizeOrigin(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "null") {
		return "", false
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return "", false
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" || parsed.User != nil {
		return "", false
	}
	return originFromURL(parsed)
}

func originFromURL(parsed *url.URL) (string, bool) {
	if parsed == nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", false
	}
	switch parsed.Scheme {
	case "http", "https", "chrome-extension":
	default:
		return "", false
	}
	return parsed.Scheme + "://" + parsed.Host, true
}
