package handler

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/karimiku/job-hunting-saas/internal/devsession"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

const (
	defaultDevAuthEmail = "dev@example.com"
	defaultDevAuthName  = "Dev User"
)

// DevAuthHandler exposes local-development-only login endpoints.
type DevAuthHandler struct {
	authenticate *useruc.Authenticate
	cfg          AuthConfig
	secret       string
}

// NewDevAuthHandler returns a local-development auth handler.
func NewDevAuthHandler(uc *useruc.Authenticate, cfg AuthConfig, secret string) *DevAuthHandler {
	return &DevAuthHandler{authenticate: uc, cfg: cfg, secret: secret}
}

// PublicRoutes registers development auth routes.
func (h *DevAuthHandler) PublicRoutes(r chi.Router) {
	r.Post("/dev/session", h.CreateSession)
}

type createDevSessionRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateSession creates a signed development session cookie without Google OAuth.
func (h *DevAuthHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	if !isLocalDevRequest(r) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var body createDevSessionRequest
	if !decodeJSONBody(w, r, &body, maxDefaultJSONBodyBytes) {
		return
	}
	email := strings.TrimSpace(body.Email)
	if email == "" {
		email = defaultDevAuthEmail
	}
	name := strings.TrimSpace(body.Name)
	if name == "" {
		name = defaultDevAuthName
	}

	out, err := h.authenticate.Execute(r.Context(), useruc.AuthenticateInput{
		Provider: "google",
		Subject:  "dev:" + email,
		Email:    email,
		Name:     name,
	})
	if err != nil {
		log.Printf("dev auth: Authenticate failed: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    devsession.Sign(out.User.ID(), h.secret),
		Path:     "/",
		Domain:   h.cfg.CookieDomain,
		MaxAge:   int(sessionMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: cookieSameSite(h.cfg),
	})

	writeJSON(w, http.StatusOK, authUserResponse{
		ID:    out.User.ID().String(),
		Email: out.User.Email().String(),
		Name:  out.User.Name().String(),
	})
}

// isLocalDevRequest reports whether the request targets a loopback host.
//
// Only r.Host (the server-observed request target) is trusted here. Origin and
// Referer are attacker-controlled request headers: an attacker can send
// Origin: http://localhost:3000 to any public host, so trusting them would let
// a production deploy issue arbitrary-user dev sessions. r.Host is the host the
// client connected to and cannot be forged into a loopback value on a public
// deployment. This is only one layer; DEV_AUTH_ENABLED, isProductionRuntime(),
// and the route-level OriginGuard remain in place (see main.go).
func isLocalDevRequest(r *http.Request) bool {
	return r.Host != "" && isLocalHost(r.Host)
}

func isLocalHost(raw string) bool {
	host := raw
	if parsedHost, _, err := net.SplitHostPort(raw); err == nil {
		host = parsedHost
	} else if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return false
		}
		host = u.Hostname()
	}
	host = strings.ToLower(strings.Trim(host, "[]"))
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
