package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type DeleteInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

// Delete は指定IDのエントリーを削除するUseCase。
type Delete struct {
	entryRepo repository.EntryRepository
}

func NewDelete(entryRepo repository.EntryRepository) *Delete {
	return &Delete{entryRepo: entryRepo}
}

// Execute はユーザーに紐づくエントリーをIDで削除する。
func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.entryRepo.Delete(ctx, input.UserID, input.EntryID)
}
