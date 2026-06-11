package repository

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// AIAccessTokenRepository は AI / MCP 連携用アクセストークンの永続化を抽象化する。
type AIAccessTokenRepository interface {
	Create(ctx context.Context, token *entity.AIAccessToken) error
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.AIAccessToken, error)
	FindActiveByHash(ctx context.Context, hash value.AIAccessTokenHash) (*entity.AIAccessToken, error)
	Revoke(ctx context.Context, userID entity.UserID, id entity.AIAccessTokenID, revokedAt time.Time) error
	TouchLastUsed(ctx context.Context, id entity.AIAccessTokenID, usedAt time.Time) error
}
