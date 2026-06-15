// Package handler は HTTP リクエストとユースケース層の橋渡しを行う。
// oapi-codegen が生成する ServerInterface を実装し、HTTP ↔ UseCase 入出力の変換のみを責務とする。
package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	useruc "github.com/karimiku/job-hunting-saas/internal/usecase/user"
)

const (
	// sessionCookieName は Session Cookie の名前。
	// Firebase Hosting の `__session` とは無関係（ここは独自ホスティング前提）。
	sessionCookieName = "session"

	// sessionMaxAge は Firebase Session Cookie の有効期間（最大14日）。
	sessionMaxAge = 14 * 24 * time.Hour

	// idTokenFreshness は受け入れる ID Token の最大鮮度。
	// 5分以内に発行されたものだけをセッション化して Session Fixation を防ぐ。
	idTokenFreshness = 5 * time.Minute
)

// AuthConfig は Cookie 発行のランタイム設定。
type AuthConfig struct {
	CookieDomain   string // 通常は空（リクエスト host に合わせる）
	CookieSecure   bool   // 本番 HTTPS では true
	CookieSameSite http.SameSite
}

// IDTokenClaims は ID Token から取り出した認証クレーム。
// Firebase 等の外部 IdP 固有型を handler 層から切り離すための DTO。
type IDTokenClaims struct {
	UID      string
	Email    string
	Name     string
	AuthTime time.Time
}

// FirebaseSessionCreator は ID Token 検証と Session Cookie 発行に必要な
// 認証バックエンドの最小インターフェース。テスト時のモック差し替え点。
// 戻り値は IdP 非依存の DTO とし、SDK 型に handler が引きずられないようにする。
type FirebaseSessionCreator interface {
	VerifyIDToken(ctx context.Context, idToken string) (*IDTokenClaims, error)
	SessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error)
}

// AuthHandler は認証関連の HTTP リクエストを受ける handler。
type AuthHandler struct {
	firebaseAuth FirebaseSessionCreator
	authenticate *useruc.Authenticate
	userRepo     repository.UserRepository
	cfg          AuthConfig
}

// NewAuthHandler は AuthHandler に必要な依存を DI して新しい AuthHandler を返す。
func NewAuthHandler(fb FirebaseSessionCreator, uc *useruc.Authenticate, userRepo repository.UserRepository, cfg AuthConfig) *AuthHandler {
	return &AuthHandler{
		firebaseAuth: fb,
		authenticate: uc,
		userRepo:     userRepo,
		cfg:          cfg,
	}
}

// PublicRoutes は認証不要なルート（ログイン/ログアウト）を登録する。
func (h *AuthHandler) PublicRoutes(r chi.Router) {
	r.Post("/auth/session", h.CreateSession)
	r.Delete("/auth/session", h.DeleteSession)
}

// ProtectedRoutes は認証必須なルート（現在のユーザー取得）を登録する。
func (h *AuthHandler) ProtectedRoutes(r chi.Router) {
	r.Get("/auth/me", h.Me)
}

type createSessionRequest struct {
	IDToken string `json:"idToken"`
}

type authUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type serverTimingMetric struct {
	name     string
	duration time.Duration
}

func appendServerTimingMetric(metrics []serverTimingMetric, name string, startedAt time.Time) []serverTimingMetric {
	return append(metrics, serverTimingMetric{name: name, duration: time.Since(startedAt)})
}

func setServerTimingHeader(w http.ResponseWriter, metrics []serverTimingMetric) {
	if len(metrics) == 0 {
		return
	}

	parts := make([]string, 0, len(metrics))
	for _, metric := range metrics {
		parts = append(parts, fmt.Sprintf("%s;dur=%.1f", metric.name, float64(metric.duration.Microseconds())/1000))
	}
	w.Header().Set("Server-Timing", strings.Join(parts, ", "))
}

// CreateSession は Firebase ID Token を検証し、User を upsert して Session Cookie を発行する。
func (h *AuthHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var body createSessionRequest
	if !decodeJSONBody(w, r, &body, maxDefaultJSONBodyBytes) {
		return
	}
	if body.IDToken == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	requestStartedAt := time.Now()
	var timingMetrics []serverTimingMetric

	stepStartedAt := time.Now()
	claims, err := h.firebaseAuth.VerifyIDToken(ctx, body.IDToken)
	timingMetrics = appendServerTimingMetric(timingMetrics, "firebase_verify_id_token", stepStartedAt)
	if err != nil {
		setServerTimingHeader(w, timingMetrics)
		log.Printf("auth: VerifyIDToken failed: %v", err)
		http.Error(w, "invalid id token", http.StatusUnauthorized)
		return
	}

	// Firebase 推奨: 鮮度チェック。古いトークンでのセッション作成を拒否する。
	if time.Since(claims.AuthTime) > idTokenFreshness {
		setServerTimingHeader(w, timingMetrics)
		http.Error(w, "recent sign-in required", http.StatusUnauthorized)
		return
	}

	name := claims.Name
	if name == "" {
		// Google アカウントで displayName が空の場合に備えたフォールバック
		name = claims.Email
	}

	stepStartedAt = time.Now()
	out, err := h.authenticate.Execute(ctx, useruc.AuthenticateInput{
		Provider: "google",
		Subject:  claims.UID,
		Email:    claims.Email,
		Name:     name,
	})
	timingMetrics = appendServerTimingMetric(timingMetrics, "user_authenticate", stepStartedAt)
	if err != nil {
		setServerTimingHeader(w, timingMetrics)
		log.Printf("auth: Authenticate failed: %v", err)
		http.Error(w, "authentication failed", http.StatusInternalServerError)
		return
	}

	stepStartedAt = time.Now()
	sessionCookie, err := h.firebaseAuth.SessionCookie(ctx, body.IDToken, sessionMaxAge)
	timingMetrics = appendServerTimingMetric(timingMetrics, "firebase_session_cookie", stepStartedAt)
	if err != nil {
		setServerTimingHeader(w, timingMetrics)
		log.Printf("auth: SessionCookie failed: %v", err)
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	// #nosec G124 -- Secure is runtime-configured: true in HTTPS deployments, false for local HTTP.
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionCookie,
		Path:     "/",
		Domain:   h.cfg.CookieDomain,
		MaxAge:   int(sessionMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   h.cfg.CookieSecure,
		SameSite: h.cookieSameSite(),
	})
	timingMetrics = appendServerTimingMetric(timingMetrics, "total", requestStartedAt)
	setServerTimingHeader(w, timingMetrics)

	writeJSON(w, http.StatusOK, authUserResponse{
		ID:    out.User.ID().String(),
		Email: out.User.Email().String(),
		Name:  out.User.Name().String(),
	})
}

// DeleteSession は Session Cookie を失効させる。認証不要（未ログインでも叩ける）。
func (h *AuthHandler) DeleteSession(w http.ResponseWriter, _ *http.Request) {
	clearSessionCookie(w, h.cfg)
	w.WriteHeader(http.StatusNoContent)
}

func clearSessionCookie(w http.ResponseWriter, cfg AuthConfig) {
	// #nosec G124 -- Secure is runtime-configured to match the cookie that was originally issued.
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		Domain:   cfg.CookieDomain,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.CookieSecure,
		SameSite: cookieSameSite(cfg),
	})
}

func (h *AuthHandler) cookieSameSite() http.SameSite {
	return cookieSameSite(h.cfg)
}

func cookieSameSite(cfg AuthConfig) http.SameSite {
	if cfg.CookieSameSite == 0 {
		return http.SameSiteLaxMode
	}
	return cfg.CookieSameSite
}

// Me は context に載った userID で User を引いて返す。
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID.IsZero() {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		log.Printf("auth: FindByID failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, authUserResponse{
		ID:    user.ID().String(),
		Email: user.Email().String(),
		Name:  user.Name().String(),
	})
}
