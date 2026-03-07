package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type UpdateInput struct {
	UserID     entity.UserID
	EntryID    entity.EntryID
	Source     string
	Status     string
	StageKind  string
	StageLabel string
	Memo       string
}

type UpdateOutput struct {
	Entry *entity.Entry
}

type Update struct {
	entryRepo repository.EntryRepository
}

func NewUpdate(entryRepo repository.EntryRepository) *Update {
	return &Update{entryRepo: entryRepo}
}

func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	validatedSource, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}

	validatedStatus, err := value.NewEntryStatus(input.Status)
	if err != nil {
		return nil, err
	}

	validatedStageKind, err := value.NewStageKind(input.StageKind)
	if err != nil {
		return nil, err
	}

	validatedStage, err := value.NewStage(validatedStageKind, input.StageLabel)
	if err != nil {
		return nil, err
	}

	entry, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	entry.UpdateSource(validatedSource)
	entry.UpdateStatus(validatedStatus)
	entry.UpdateStage(validatedStage)
	entry.UpdateMemo(input.Memo)

	if err := uc.entryRepo.Save(ctx, entry); err != nil {
		return nil, err
	}

	return &UpdateOutput{Entry: entry}, nil
}
