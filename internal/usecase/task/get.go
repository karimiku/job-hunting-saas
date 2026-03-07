package task

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetInput struct {
	UserID entity.UserID
	TaskID entity.TaskID
}

type GetOutput struct {
	Task *entity.Task
}

type Get struct {
	taskRepo repository.TaskRepository
}

func NewGet(taskRepo repository.TaskRepository) *Get {
	return &Get{taskRepo: taskRepo}
}

func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	task, err := uc.taskRepo.FindByID(ctx, input.UserID, input.TaskID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{Task: task}, nil
}
