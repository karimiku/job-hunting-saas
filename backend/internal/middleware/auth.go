// Package middleware は HTTP リクエストの前処理 (認証等) を担う。
package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// contextKey は独自型を使い、他パッケージとのcontextキー衝突を防ぐ。
type contextKey string

const userIDKey contextKey = "userID"

// SessionCookieName は Auth ミドルウェアと AuthHandler で共有する Cookie 名。
const SessionCookieName = "session"

// SetUserID は認証済みユーザーのIDをcontextに埋め込む。
func SetUserID(ctx context.Context, userID entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID はcontextから認証済みユーザーのIDを取り出す。
// 未認証の場合はゼロ値を返す（IsZero() で判定可）。
func GetUserID(ctx context.Context) entity.UserID {
	userID, found := ctx.Value(userIDKey).(entity.UserID)
	if !found {
		return entity.UserID{}
	}
	return userID
}

// SessionClaims は Session Cookie から取り出した認証クレーム。
// Firebase 等の外部 IdP 固有型を middleware 層から切り離すための DTO。
type SessionClaims struct {
	UID string
}

// FirebaseSessionVerifier は Session Cookie 検証に必要な認証バックエンドの最小インターフェース。
// テスト時のモック差し替え点。戻り値は IdP 非依存の DTO。
type FirebaseSessionVerifier interface {
	VerifySessionCookie(ctx context.Context, sessionCookie string) (*SessionClaims, error)
}

// BearerTokenVerifier は Authorization: Bearer で渡されたAI連携トークンを userID に解決する。
type BearerTokenVerifier interface {
	VerifyBearerToken(ctx context.Context, rawToken string) (entity.UserID, error)
}

// NewAuth は Session Cookie を検証して userID を context に埋め込む chi ミドルウェアを返す。
//
// フロー:
//  1. Cookie 取得 → なければ 401
//  2. Firebase Admin SDK で Session Cookie 検証 → 失敗は 401
//  3. Firebase UID → external_identities → users の順で UserID を解決
//  4. context に UserID をセットして次へ
func NewAuth(fbAuth FirebaseSessionVerifier, extIDRepo repository.ExternalIdentityRepository) func(http.Handler) http.Handler {
	return NewAuthWithBearer(fbAuth, extIDRepo, nil)
}

// NewAuthWithBearer は Session Cookie に加えて Authorization: Bearer も受け付ける。
//
// Browser は従来通り Session Cookie、Claude/Codex/MCP 等の外部AIクライアントは
// AI連携トークンを Bearer で送る。どちらも最終的に context の UserID に正規化する。
func NewAuthWithBearer(
	fbAuth FirebaseSessionVerifier,
	extIDRepo repository.ExternalIdentityRepository,
	bearerVerifier BearerTokenVerifier,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if header := strings.TrimSpace(r.Header.Get("Authorization")); header != "" {
				userID, ok, err := verifyBearer(ctx, header, bearerVerifier)
				if err != nil {
					if errors.Is(err, value.ErrAIAccessTokenInvalid) {
						http.Error(w, "unauthenticated", http.StatusUnauthorized)
						return
					}
					log.Printf("auth middleware: VerifyBearerToken: %v", err)
					http.Error(w, "internal error", http.StatusInternalServerError)
					return
				}
				if !ok {
					http.Error(w, "unauthenticated", http.StatusUnauthorized)
					return
				}
				ctx = SetUserID(ctx, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			claims, err := fbAuth.VerifySessionCookie(ctx, cookie.Value)
			if err != nil {
				// 失効・改ざん・期限切れを区別せずに 401 を返す
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			identity, err := extIDRepo.FindByProviderAndSubject(ctx, value.AuthProviderGoogle(), claims.UID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					// Session は有効だが DB にユーザーがいない異常系
					http.Error(w, "unauthenticated", http.StatusUnauthorized)
					return
				}
				log.Printf("auth middleware: FindByProviderAndSubject: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			ctx = SetUserID(ctx, identity.UserID())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyBearer(ctx context.Context, header string, verifier BearerTokenVerifier) (entity.UserID, bool, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(header, bearerPrefix) {
		return entity.UserID{}, false, nil
	}
	rawToken := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
	if rawToken == "" || verifier == nil {
		return entity.UserID{}, false, nil
	}
	userID, err := verifier.VerifyBearerToken(ctx, rawToken)
	if err != nil {
		return entity.UserID{}, false, err
	}
	return userID, true, nil
}
