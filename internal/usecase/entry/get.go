package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

type GetOutput struct {
	Entry *entity.Entry
}

type Get struct {
	entryRepo repository.EntryRepository
}

func NewGet(entryRepo repository.EntryRepository) *Get {
	return &Get{entryRepo: entryRepo}
}

func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	entry, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{Entry: entry}, nil
}
