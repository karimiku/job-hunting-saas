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

type Delete struct {
	taskRepo repository.TaskRepository
}

func NewDelete(taskRepo repository.TaskRepository) *Delete {
	return &Delete{taskRepo: taskRepo}
}

func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.taskRepo.Delete(ctx, input.UserID, input.TaskID)
}
