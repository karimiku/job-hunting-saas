package aiaccesstoken

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListInput は AI 連携トークン一覧への入力。
type ListInput struct {
	UserID entity.UserID
}

// ListOutput は AI 連携トークン一覧の出力。
type ListOutput struct {
	Tokens []*entity.AIAccessToken
}

// List はユーザーが発行した AI 連携トークンを一覧するユースケース。
type List struct {
	repo repository.AIAccessTokenRepository
}

// NewList は List ユースケースを生成する。
func NewList(repo repository.AIAccessTokenRepository) *List {
	return &List{repo: repo}
}

// Execute は userID 所有のトークンを作成日時の新しい順で返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	if input.UserID.IsZero() {
		return nil, ErrUserIDRequired
	}
	tokens, err := uc.repo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	return &ListOutput{Tokens: tokens}, nil
}
