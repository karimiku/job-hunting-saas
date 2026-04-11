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

// Update は既存エントリーのSource・Status・Stage・Memoを更新するUseCase。
type Update struct {
	entryRepo repository.EntryRepository
}

func NewUpdate(entryRepo repository.EntryRepository) *Update {
	return &Update{entryRepo: entryRepo}
}

// Execute は各値をバリデーションし、既存Entryを取得して更新する。
func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	source, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}

	status, err := value.NewEntryStatus(input.Status)
	if err != nil {
		return nil, err
	}

	stageKind, err := value.NewStageKind(input.StageKind)
	if err != nil {
		return nil, err
	}

	stage, err := value.NewStage(stageKind, input.StageLabel)
	if err != nil {
		return nil, err
	}

	e, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	e.UpdateSource(source)
	e.UpdateStatus(status)
	e.UpdateStage(stage)
	e.UpdateMemo(input.Memo)

	if err := uc.entryRepo.Save(ctx, e); err != nil {
		return nil, err
	}

	return &UpdateOutput{Entry: e}, nil
}
