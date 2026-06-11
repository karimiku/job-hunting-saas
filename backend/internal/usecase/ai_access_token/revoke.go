package aiaccesstoken

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// RevokeInput は AI 連携トークン失効への入力。
type RevokeInput struct {
	UserID  entity.UserID
	TokenID entity.AIAccessTokenID
}

// Revoke は AI 連携トークンを失効するユースケース。
type Revoke struct {
	repo repository.AIAccessTokenRepository
}

// NewRevoke は Revoke ユースケースを生成する。
func NewRevoke(repo repository.AIAccessTokenRepository) *Revoke {
	return &Revoke{repo: repo}
}

// Execute は所有者だけがトークンを失効できるよう userID で絞り込む。
func (uc *Revoke) Execute(ctx context.Context, input RevokeInput) error {
	if input.UserID.IsZero() {
		return ErrUserIDRequired
	}
	return uc.repo.Revoke(ctx, input.UserID, input.TokenID, time.Now())
}
