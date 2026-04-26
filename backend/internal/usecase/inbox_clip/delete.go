package inboxclip

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// DeleteInput は Delete ユースケースへの入力。
type DeleteInput struct {
	UserID entity.UserID
	ClipID entity.InboxClipID
}

// Delete はクリップを削除するユースケース。
type Delete struct {
	repo repository.InboxClipRepository
}

// NewDelete は Delete ユースケースを生成する。
func NewDelete(repo repository.InboxClipRepository) *Delete {
	return &Delete{repo: repo}
}

// Execute は所有権を確認してクリップを削除する。
func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.repo.Delete(ctx, input.UserID, input.ClipID)
}
