package task

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
	Tasks []*entity.Task
}

type List struct {
	taskRepo repository.TaskRepository
}

func NewList(taskRepo repository.TaskRepository) *List {
	return &List{taskRepo: taskRepo}
}

func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	tasks, err := uc.taskRepo.ListByEntryID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Tasks: tasks}, nil
}
