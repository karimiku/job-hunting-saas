package task

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
	Title   string
	Type    string
	DueDate *time.Time
	Memo    string
}

type CreateOutput struct {
	Task *entity.Task
}

type Create struct {
	taskRepo  repository.TaskRepository
	entryRepo repository.EntryRepository
}

func NewCreate(taskRepo repository.TaskRepository, entryRepo repository.EntryRepository) *Create {
	return &Create{taskRepo: taskRepo, entryRepo: entryRepo}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	// 指定されたEntryが存在し、かつ操作ユーザーが所有していることを検証する
	if _, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID); err != nil {
		return nil, err
	}

	validatedTitle, err := value.NewTaskTitle(input.Title)
	if err != nil {
		return nil, err
	}

	validatedType, err := value.NewTaskType(input.Type)
	if err != nil {
		return nil, err
	}

	task := entity.NewTask(input.EntryID, validatedTitle, validatedType)

	if input.DueDate != nil {
		task.SetDueDate(*input.DueDate)
	}

	if input.Memo != "" {
		task.UpdateMemo(input.Memo)
	}

	if err := uc.taskRepo.Save(ctx, task); err != nil {
		return nil, err
	}

	return &CreateOutput{Task: task}, nil
}
