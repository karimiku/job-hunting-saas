package selectionflow

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// GetInput は選考フロー取得ユースケースへの入力。
type GetInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

// GetOutput は選考フロー取得ユースケースの出力。
type GetOutput struct {
	SelectionFlow *entity.SelectionFlow
}

// Get はEntryに紐づく選考フローを取得する。
type Get struct {
	flowRepo  repository.SelectionFlowRepository
	entryRepo repository.EntryRepository
}

// NewGet は Get ユースケースを生成する。
func NewGet(flowRepo repository.SelectionFlowRepository, entryRepo repository.EntryRepository) *Get {
	return &Get{flowRepo: flowRepo, entryRepo: entryRepo}
}

// Execute はEntry所有権を確認し、選考フローを返す。
func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}
	flow, err := uc.flowRepo.FindByEntryID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}
	return &GetOutput{SelectionFlow: flow}, nil
}
