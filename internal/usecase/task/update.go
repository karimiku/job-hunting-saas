package task

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type UpdateInput struct {
	UserID  entity.UserID
	TaskID  entity.TaskID
	Title   string
	Type    string
	Status  string
	DueDate *time.Time
	Notify  bool
	Memo    string
}

type UpdateOutput struct {
	Task *entity.Task
}

type Update struct {
	taskRepo repository.TaskRepository
}

func NewUpdate(taskRepo repository.TaskRepository) *Update {
	return &Update{taskRepo: taskRepo}
}

func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	taskTitle, err := value.NewTaskTitle(input.Title)
	if err != nil {
		return nil, err
	}

	taskType, err := value.NewTaskType(input.Type)
	if err != nil {
		return nil, err
	}

	validatedStatus, err := value.NewTaskStatus(input.Status)
	if err != nil {
		return nil, err
	}

	task, err := uc.taskRepo.FindByID(ctx, input.UserID, input.TaskID)
	if err != nil {
		return nil, err
	}

	task.UpdateTitle(taskTitle)
	task.UpdateTaskType(taskType)

	if validatedStatus.IsDone() {
		task.Complete()
	} else {
		task.Uncomplete()
	}

	if input.DueDate != nil {
		task.SetDueDate(*input.DueDate)
	} else {
		task.ClearDueDate()
	}

	task.SetNotify(input.Notify)
	task.UpdateMemo(input.Memo)

	if err := uc.taskRepo.Save(ctx, task); err != nil {
		return nil, err
	}

	return &UpdateOutput{Task: task}, nil
}
