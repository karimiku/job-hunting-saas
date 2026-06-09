package task

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListAllInput はユーザーの全 Task 一覧ユースケースへの入力。
type ListAllInput struct {
	UserID entity.UserID
}

// ListAllOutput はユーザーの全 Task 一覧ユースケースの出力。
type ListAllOutput struct {
	Tasks []*entity.Task
}

// ListAll はユーザー所有の全 Task を取得する UseCase。
type ListAll struct {
	taskRepo repository.TaskRepository
}

// NewListAll は ListAll ユースケースを生成する。
func NewListAll(taskRepo repository.TaskRepository) *ListAll {
	return &ListAll{taskRepo: taskRepo}
}

// Execute はユーザー所有の全 Task を検索して返す。
func (uc *ListAll) Execute(ctx context.Context, input ListAllInput) (*ListAllOutput, error) {
	tasks, err := uc.taskRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &ListAllOutput{Tasks: tasks}, nil
}
