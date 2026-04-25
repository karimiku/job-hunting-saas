// Package task はタスクに対するユースケース群を提供する。
package task

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateInput は TaskCreate ユースケースへの入力。
type CreateInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
	Title   string
	Type    string
	DueDate *time.Time
	Memo    string
}

// CreateOutput は TaskCreate ユースケースの出力。
type CreateOutput struct {
	Task *entity.Task
}

// Create は新しいタスクを登録するUseCase。
type Create struct {
	taskRepo  repository.TaskRepository
	entryRepo repository.EntryRepository
}

// NewCreate は TaskCreate ユースケースを生成する。
func NewCreate(taskRepo repository.TaskRepository, entryRepo repository.EntryRepository) *Create {
	return &Create{taskRepo: taskRepo, entryRepo: entryRepo}
}

// Execute はEntryIDの存在・所有を検証し、Title/TaskTypeをバリデーションして新規Taskを生成・永続化する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}

	title, err := value.NewTaskTitle(input.Title)
	if err != nil {
		return nil, err
	}

	taskType, err := value.NewTaskType(input.Type)
	if err != nil {
		return nil, err
	}

	t := entity.NewTask(input.EntryID, title, taskType)

	if input.DueDate != nil {
		t.SetDueDate(*input.DueDate)
	}

	if input.Memo != "" {
		t.UpdateMemo(input.Memo)
	}

	if err := uc.taskRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return &CreateOutput{Task: t}, nil
}
