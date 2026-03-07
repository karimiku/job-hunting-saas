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

type Delete struct {
	entryRepo repository.EntryRepository
}

func NewDelete(entryRepo repository.EntryRepository) *Delete {
	return &Delete{entryRepo: entryRepo}
}

func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.entryRepo.Delete(ctx, input.UserID, input.EntryID)
}
