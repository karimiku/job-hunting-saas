package repository

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// AIAccessTokenRepository はAIクライアント連携用アクセストークンの永続化を抽象化する。
type AIAccessTokenRepository interface {
	Save(ctx context.Context, token *entity.AIAccessToken) error
	FindActiveByHash(ctx context.Context, tokenHash string) (*entity.AIAccessToken, error)
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.AIAccessToken, error)
	TouchLastUsed(ctx context.Context, id entity.AIAccessTokenID, usedAt time.Time) error
	Revoke(ctx context.Context, userID entity.UserID, id entity.AIAccessTokenID, revokedAt time.Time) error
}
