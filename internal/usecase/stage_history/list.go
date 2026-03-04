package stagehistory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type ListInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

type ListOutput struct {
	StageHistories []*entity.StageHistory
}

// List はエントリーに紐づく選考履歴一覧を取得するUseCase。
type List struct {
	historyRepo repository.StageHistoryRepository
	entryRepo   repository.EntryRepository
}

func NewList(historyRepo repository.StageHistoryRepository, entryRepo repository.EntryRepository) *List {
	return &List{historyRepo: historyRepo, entryRepo: entryRepo}
}

// Execute はユーザーのエントリー所有を検証し、選考履歴一覧を返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}

	histories, err := uc.historyRepo.ListByEntryID(ctx, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{StageHistories: histories}, nil
}
