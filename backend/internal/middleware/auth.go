// Package middleware は HTTP リクエストの前処理 (認証等) を担う。
package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/devsession"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// contextKey は独自型を使い、他パッケージとのcontextキー衝突を防ぐ。
type contextKey string

const userIDKey contextKey = "userID"
const authMethodKey contextKey = "authMethod"

// SessionCookieName は Auth ミドルウェアと AuthHandler で共有する Cookie 名。
const SessionCookieName = "session"

// AuthMethod は認証済みリクエストがどの方式で認証されたかを表す。
type AuthMethod string

const (
	// AuthMethodUnknown は認証方式が context に無い状態。
	AuthMethodUnknown AuthMethod = ""
	// AuthMethodSession はブラウザ向け Session Cookie 認証。
	AuthMethodSession AuthMethod = "session"
	// AuthMethodBearer は AI / MCP 連携用 Bearer token 認証。
	AuthMethodBearer AuthMethod = "bearer"
)

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

// SetAuthMethod は認証方式をcontextに埋め込む。
func SetAuthMethod(ctx context.Context, method AuthMethod) context.Context {
	return context.WithValue(ctx, authMethodKey, method)
}

// GetAuthMethod はcontextから認証方式を取り出す。
func GetAuthMethod(ctx context.Context) AuthMethod {
	method, found := ctx.Value(authMethodKey).(AuthMethod)
	if !found {
		return AuthMethodUnknown
	}
	return method
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

// BearerTokenVerifier は Authorization: Bearer で渡されたトークンを userID に解決する。
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
	return NewAuthWithBearerAndDevSession(fbAuth, extIDRepo, bearerVerifier, "")
}

// NewAuthWithBearerAndDevSession は開発専用の署名付き session cookie も受け付ける。
// devSessionSecret が空なら通常の Firebase Session Cookie 検証だけを行う。
func NewAuthWithBearerAndDevSession(
	fbAuth FirebaseSessionVerifier,
	extIDRepo repository.ExternalIdentityRepository,
	bearerVerifier BearerTokenVerifier,
	devSessionSecret string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if header := strings.TrimSpace(r.Header.Get("Authorization")); header != "" {
				startedAt := time.Now()
				userID, ok, err := verifyBearer(ctx, header, bearerVerifier)
				addServerTimingSince(ctx, "auth_bearer", startedAt)
				if err != nil {
					if isInvalidBearerTokenError(err) {
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
				ctx = SetAuthMethod(SetUserID(ctx, userID), AuthMethodBearer)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}
			if fbAuth == nil {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			if userID, ok := devsession.Verify(cookie.Value, devSessionSecret); ok {
				ctx = SetAuthMethod(SetUserID(ctx, userID), AuthMethodSession)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			if fbAuth == nil {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			startedAt := time.Now()
			claims, err := fbAuth.VerifySessionCookie(ctx, cookie.Value)
			addServerTimingSince(ctx, "firebase_verify_session_cookie", startedAt)
			if err != nil {
				// 失効・改ざん・期限切れを区別せずに 401 を返す
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}

			startedAt = time.Now()
			identity, err := extIDRepo.FindByProviderAndSubject(ctx, value.AuthProviderGoogle(), claims.UID)
			addServerTimingSince(ctx, "external_identity_lookup", startedAt)
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

			ctx = SetAuthMethod(SetUserID(ctx, identity.UserID()), AuthMethodSession)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// NewChainedBearerTokenVerifier は複数の Bearer verifier を順番に試す。
func NewChainedBearerTokenVerifier(verifiers ...BearerTokenVerifier) BearerTokenVerifier {
	filtered := make([]BearerTokenVerifier, 0, len(verifiers))
	for _, verifier := range verifiers {
		if verifier != nil {
			filtered = append(filtered, verifier)
		}
	}
	return chainedBearerTokenVerifier{verifiers: filtered}
}

type chainedBearerTokenVerifier struct {
	verifiers []BearerTokenVerifier
}

func (v chainedBearerTokenVerifier) VerifyBearerToken(ctx context.Context, rawToken string) (entity.UserID, error) {
	for _, verifier := range v.verifiers {
		userID, err := verifier.VerifyBearerToken(ctx, rawToken)
		if err == nil {
			return userID, nil
		}
		if isInvalidBearerTokenError(err) {
			continue
		}
		return entity.UserID{}, err
	}
	return entity.UserID{}, value.ErrAuthTokenInvalid
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

func isInvalidBearerTokenError(err error) bool {
	return errors.Is(err, value.ErrAuthTokenInvalid) || errors.Is(err, value.ErrAIAccessTokenInvalid)
}
