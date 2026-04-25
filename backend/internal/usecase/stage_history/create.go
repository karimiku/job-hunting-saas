// Package stagehistory は選考ステージ履歴のユースケース群を提供する。
package stagehistory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateInput は StageHistoryCreate ユースケースへの入力。
type CreateInput struct {
	UserID    entity.UserID
	EntryID   entity.EntryID
	StageKind string
	Label     string
	Note      string
}

// CreateOutput は StageHistoryCreate ユースケースの出力。
type CreateOutput struct {
	StageHistory *entity.StageHistory
}

// Create は選考フェーズ履歴を追加するUseCase。
type Create struct {
	historyRepo repository.StageHistoryRepository
	entryRepo   repository.EntryRepository
}

// NewCreate は StageHistoryCreate ユースケースを生成する。
func NewCreate(historyRepo repository.StageHistoryRepository, entryRepo repository.EntryRepository) *Create {
	return &Create{historyRepo: historyRepo, entryRepo: entryRepo}
}

// Execute はEntryIDの存在・所有を検証し、Stage値オブジェクトをバリデーションしてStageHistoryを生成・永続化する。
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

	history := entity.NewStageHistory(input.EntryID, stage, input.Note)

	if err := uc.historyRepo.Create(ctx, history); err != nil {
		return nil, err
	}

	return &CreateOutput{StageHistory: history}, nil
}
