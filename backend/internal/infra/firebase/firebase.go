// Package firebase は Firebase Admin SDK の初期化と、handler / middleware が
// 期待する DTO 型へのアダプタを提供する。
// SDK 固有型 (*auth.Token 等) はこの層に閉じ込め、上位層は IdP 非依存の DTO だけを扱う。
package firebase

import (
	"context"
	"fmt"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"

	"github.com/karimiku/job-hunting-saas/internal/handler"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
)

// Client は Firebase Admin SDK の Auth クライアントを保持する。
type Client struct {
	Auth *auth.Client
}

// NewClient は Firebase Admin SDK を初期化して Client を返す。
//
// credentialsPath が空の場合は ADC（Application Default Credentials）を利用する。
// 通常は GOOGLE_APPLICATION_CREDENTIALS env か、ここに service account JSON パスを渡す。
func NewClient(ctx context.Context, credentialsPath, projectID string) (*Client, error) {
	cfg := &firebase.Config{ProjectID: projectID}

	var opts []option.ClientOption
	if credentialsPath != "" {
		// option.WithCredentialsFile は遅延読み込みかつ TOCTOU 系の懸念から deprecated 扱い。
		// 起動時に明示的に読んで JSON バイト列で渡す。
		// credentialsPath はオペレータ管理の env var (FIREBASE_CREDENTIALS_FILE) 経由でのみ
		// 入る値であり、外部入力ではないため G304 は false positive として抑止する。
		data, err := os.ReadFile(credentialsPath) // #nosec G304
		if err != nil {
			return nil, fmt.Errorf("firebase: read credentials %q: %w", credentialsPath, err)
		}
		// WithCredentialsJSON も google.golang.org/api/option では deprecated 扱いだが、
		// 代替の golang.org/x/oauth2/google.CredentialsFromJSON 経由は scope 指定が必要で
		// firebase.NewApp の挙動と完全互換にするには一手間かかる。
		// 現状のままでも実行時の挙動は変わらないので staticcheck SA1019 を抑止する。
		//nolint:staticcheck // SA1019: alternative path needs explicit scope wiring; revisit when Firebase SDK guidance updates.
		opts = append(opts, option.WithCredentialsJSON(data))
	}

	app, err := firebase.NewApp(ctx, cfg, opts...)
	if err != nil {
		return nil, fmt.Errorf("firebase: init app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("firebase: init auth client: %w", err)
	}
	return &Client{Auth: authClient}, nil
}

// SessionCreator は handler.FirebaseSessionCreator を実装する Firebase アダプタ。
// SDK の *auth.Token を handler.IDTokenClaims (DTO) に変換して返す。
type SessionCreator struct {
	client *auth.Client
}

// NewSessionCreator は SessionCreator を生成する。
func NewSessionCreator(client *auth.Client) *SessionCreator {
	return &SessionCreator{client: client}
}

// VerifyIDToken は ID Token を検証し、必要なクレームだけを DTO に詰めて返す。
func (a *SessionCreator) VerifyIDToken(ctx context.Context, idToken string) (*handler.IDTokenClaims, error) {
	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)
	return &handler.IDTokenClaims{
		UID:      token.UID,
		Email:    email,
		Name:     name,
		AuthTime: time.Unix(token.AuthTime, 0),
	}, nil
}

// SessionCookie は Firebase Admin SDK の SessionCookie 発行をそのまま委譲する。
func (a *SessionCreator) SessionCookie(ctx context.Context, idToken string, expiresIn time.Duration) (string, error) {
	return a.client.SessionCookie(ctx, idToken, expiresIn)
}

// SessionVerifier は middleware.FirebaseSessionVerifier を実装する Firebase アダプタ。
type SessionVerifier struct {
	client *auth.Client
}

// NewSessionVerifier は SessionVerifier を生成する。
func NewSessionVerifier(client *auth.Client) *SessionVerifier {
	return &SessionVerifier{client: client}
}

// VerifySessionCookie は Session Cookie を検証し、必要なクレームだけを DTO に詰めて返す。
func (a *SessionVerifier) VerifySessionCookie(ctx context.Context, sessionCookie string) (*middleware.SessionClaims, error) {
	token, err := a.client.VerifySessionCookie(ctx, sessionCookie)
	if err != nil {
		return nil, err
	}
	return &middleware.SessionClaims{UID: token.UID}, nil
}
