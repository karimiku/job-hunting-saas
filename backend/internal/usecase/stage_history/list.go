package stagehistory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListInput は StageHistoryList ユースケースへの入力。
type ListInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

// ListOutput は StageHistoryList ユースケースの出力。
type ListOutput struct {
	StageHistories []*entity.StageHistory
}

// List はエントリーに紐づく選考フェーズ履歴一覧を取得するUseCase。
type List struct {
	historyRepo repository.StageHistoryRepository
	entryRepo   repository.EntryRepository
}

// NewList は StageHistoryList ユースケースを生成する。
func NewList(historyRepo repository.StageHistoryRepository, entryRepo repository.EntryRepository) *List {
	return &List{historyRepo: historyRepo, entryRepo: entryRepo}
}

// Execute はユーザー・エントリーに紐づく選考フェーズ履歴を返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	// EntryのuserID所有権を検証
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}

	histories, err := uc.historyRepo.ListByEntryID(ctx, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{StageHistories: histories}, nil
}
