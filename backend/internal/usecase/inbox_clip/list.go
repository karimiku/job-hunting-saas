package inboxclip

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListInput は List ユースケースへの入力。
type ListInput struct {
	UserID entity.UserID
}

// ListOutput は List ユースケースの出力。
type ListOutput struct {
	Clips []*entity.InboxClip
}

// List はユーザーのクリップ一覧を取得するユースケース。
type List struct {
	repo repository.InboxClipRepository
}

// NewList は List ユースケースを生成する。
func NewList(repo repository.InboxClipRepository) *List {
	return &List{repo: repo}
}

// Execute は userID 所有のクリップを保存日時の新しい順で返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	clips, err := uc.repo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	return &ListOutput{Clips: clips}, nil
}
