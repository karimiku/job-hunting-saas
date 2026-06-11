package esmemo

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

const (
	DefaultListLimit = 50
	MaxListLimit     = 100
)

// ListInput はESメモ一覧ユースケースへの入力。
type ListInput struct {
	UserID entity.UserID
	Limit  int32
}

// ListOutput はESメモ一覧ユースケースの出力。
type ListOutput struct {
	Memos []*entity.ESMemo
}

// List はユーザーのESメモ一覧を取得するUseCase。
type List struct {
	memoRepo repository.ESMemoRepository
}

// NewList はESメモ一覧ユースケースを生成する。
func NewList(memoRepo repository.ESMemoRepository) *List {
	return &List{memoRepo: memoRepo}
}

// Execute はESメモを新しい順に返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = DefaultListLimit
	}
	if limit > MaxListLimit {
		limit = MaxListLimit
	}
	memos, err := uc.memoRepo.ListByUserID(ctx, input.UserID, limit)
	if err != nil {
		return nil, err
	}
	return &ListOutput{Memos: memos}, nil
}
