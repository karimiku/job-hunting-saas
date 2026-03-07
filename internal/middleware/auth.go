package middleware

import (
	"context"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// contextKey は独自型を使い、他パッケージとのcontextキー衝突を防ぐ。
type contextKey string

const userIDKey contextKey = "userID"

// SetUserID は認証済みユーザーのIDをcontextに埋め込む。
// Auth ミドルウェアがトークン検証後に呼び出す。
func SetUserID(ctx context.Context, userID entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID はcontextから認証済みユーザーのIDを取り出す。
// 未認証の場合はゼロ値を返す（Auth実装後は到達しない想定）。
func GetUserID(ctx context.Context) entity.UserID {
	userID, found := ctx.Value(userIDKey).(entity.UserID)
	if !found {
		return entity.UserID{}
	}
	return userID
}

// Auth はリクエストごとにセッションの有効性を検証し、
// 認証済みユーザーのIDをcontextに埋め込むミドルウェア。
// TODO: Cookieからセッショントークンを取り出し、検証し、SetUserIDでcontextに埋め込む。
//
//	未認証の場合は 401 Unauthorized を返して後続handlerを呼ばない。
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
