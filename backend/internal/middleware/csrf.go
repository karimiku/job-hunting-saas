package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// NewOriginGuard rejects unsafe browser requests unless Origin or Referer
// resolves to one of the explicitly allowed origins.
func NewOriginGuard(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := newOriginAllowlist(allowedOrigins)
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
	allowed := newOriginAllowlist(allowedOrigins)
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

type originAllowlist struct {
	exact    map[string]struct{}
	wildcard []originWildcard
}

type originWildcard struct {
	scheme string
	suffix string
}

func newOriginAllowlist(origins []string) originAllowlist {
	allowlist := originAllowlist{
		exact: make(map[string]struct{}, len(origins)),
	}
	for _, raw := range origins {
		if wildcard, ok := normalizeOriginWildcard(raw); ok {
			allowlist.wildcard = append(allowlist.wildcard, wildcard)
			continue
		}
		if origin, ok := normalizeOrigin(raw); ok {
			allowlist.exact[origin] = struct{}{}
		}
	}
	return allowlist
}

// OriginAllowed reports whether rawOrigin matches the configured exact origins
// or wildcard host patterns such as https://*.vercel.app.
func OriginAllowed(rawOrigin string, allowedOrigins []string) bool {
	return newOriginAllowlist(allowedOrigins).allows(rawOrigin)
}

func (a originAllowlist) allows(rawOrigin string) bool {
	origin, ok := normalizeOrigin(rawOrigin)
	if !ok {
		return false
	}
	if _, found := a.exact[origin]; found {
		return true
	}

	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Host)
	for _, wildcard := range a.wildcard {
		if parsed.Scheme == wildcard.scheme &&
			strings.HasSuffix(host, wildcard.suffix) &&
			len(host) > len(strings.TrimPrefix(wildcard.suffix, ".")) {
			return true
		}
	}
	return false
}

func requestOriginAllowed(r *http.Request, allowed originAllowlist) bool {
	if len(allowed.exact) == 0 && len(allowed.wildcard) == 0 {
		return false
	}
	if rawOrigin := r.Header.Get("Origin"); rawOrigin != "" {
		return allowed.allows(rawOrigin)
	}
	if origin, ok := refererOrigin(r.Header.Get("Referer")); ok {
		return allowed.allows(origin)
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

func normalizeOriginWildcard(raw string) (originWildcard, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.EqualFold(raw, "null") {
		return originWildcard{}, false
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return originWildcard{}, false
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return originWildcard{}, false
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" || parsed.User != nil {
		return originWildcard{}, false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return originWildcard{}, false
	}
	host := strings.ToLower(parsed.Host)
	if !strings.HasPrefix(host, "*.") {
		return originWildcard{}, false
	}
	suffix := strings.TrimPrefix(host, "*")
	if suffix == "." || strings.Contains(strings.TrimPrefix(suffix, "."), "*") {
		return originWildcard{}, false
	}
	return originWildcard{scheme: parsed.Scheme, suffix: suffix}, true
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
	return parsed.Scheme + "://" + strings.ToLower(parsed.Host), true
}
