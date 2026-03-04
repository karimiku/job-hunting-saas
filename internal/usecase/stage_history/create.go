package stagehistory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateInput struct {
	UserID    entity.UserID
	EntryID   entity.EntryID
	StageKind string
	Label     string
	Note      string
}

type CreateOutput struct {
	StageHistory *entity.StageHistory
}

// Create は新しい選考履歴を登録するUseCase。
type Create struct {
	historyRepo repository.StageHistoryRepository
	entryRepo   repository.EntryRepository
}

func NewCreate(historyRepo repository.StageHistoryRepository, entryRepo repository.EntryRepository) *Create {
	return &Create{historyRepo: historyRepo, entryRepo: entryRepo}
}

// Execute はEntryIDの存在・所有を検証し、Stageをバリデーションして新規StageHistoryを生成・永続化する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}

	kind, err := value.NewStageKind(input.StageKind)
	if err != nil {
		return nil, err
	}

	stage, err := value.NewStage(kind, input.Label)
	if err != nil {
		return nil, err
	}

	h := entity.NewStageHistory(input.EntryID, stage, input.Note)

	if err := uc.historyRepo.Create(ctx, h); err != nil {
		return nil, err
	}

	return &CreateOutput{StageHistory: h}, nil
}
