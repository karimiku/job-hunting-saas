package task

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type DeleteInput struct {
	UserID entity.UserID
	TaskID entity.TaskID
}

// Delete は指定IDのタスクを削除するUseCase。
type Delete struct {
	taskRepo repository.TaskRepository
}

func NewDelete(taskRepo repository.TaskRepository) *Delete {
	return &Delete{taskRepo: taskRepo}
}

// Execute はユーザーに紐づくタスクをIDで削除する。
func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.taskRepo.Delete(ctx, input.UserID, input.TaskID)
}
