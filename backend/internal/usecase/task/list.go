package task

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListInput は TaskList ユースケースへの入力。
type ListInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

// ListOutput は TaskList ユースケースの出力。
type ListOutput struct {
	Tasks []*entity.Task
}

// List はエントリーに紐づくタスク一覧を取得するUseCase。
type List struct {
	taskRepo repository.TaskRepository
}

// NewList は TaskList ユースケースを生成する。
func NewList(taskRepo repository.TaskRepository) *List {
	return &List{taskRepo: taskRepo}
}

// Execute はユーザー・エントリーに紐づくタスク一覧を検索して返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	tasks, err := uc.taskRepo.ListByEntryID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Tasks: tasks}, nil
}
