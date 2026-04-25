package task

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// UpdateInput は TaskUpdate ユースケースへの入力。
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

// UpdateOutput は TaskUpdate ユースケースの出力。
type UpdateOutput struct {
	Task *entity.Task
}

// Update は既存タスクのTitle・Type・Status・DueDate・Notify・Memoを更新するUseCase。
type Update struct {
	taskRepo repository.TaskRepository
}

// NewUpdate は TaskUpdate ユースケースを生成する。
func NewUpdate(taskRepo repository.TaskRepository) *Update {
	return &Update{taskRepo: taskRepo}
}

// Execute は各値をバリデーションし、既存Taskを取得して更新する。
func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	title, err := value.NewTaskTitle(input.Title)
	if err != nil {
		return nil, err
	}

	taskType, err := value.NewTaskType(input.Type)
	if err != nil {
		return nil, err
	}

	status, err := value.NewTaskStatus(input.Status)
	if err != nil {
		return nil, err
	}

	t, err := uc.taskRepo.FindByID(ctx, input.UserID, input.TaskID)
	if err != nil {
		return nil, err
	}

	t.UpdateTitle(title)
	t.UpdateTaskType(taskType)

	if status.IsDone() {
		t.Complete()
	} else {
		t.Uncomplete()
	}

	if input.DueDate != nil {
		t.SetDueDate(*input.DueDate)
	} else {
		t.ClearDueDate()
	}

	t.SetNotify(input.Notify)
	t.UpdateMemo(input.Memo)

	if err := uc.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return &UpdateOutput{Task: t}, nil
}
