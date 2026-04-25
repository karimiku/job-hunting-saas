// Package firebase は Firebase Admin SDK の初期化を行う。
// 他パッケージからは *auth.Client を介して ID Token 検証や Session Cookie 発行を行う。
package firebase

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
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
