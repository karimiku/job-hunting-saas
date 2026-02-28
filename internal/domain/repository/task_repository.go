package repository

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// TaskRepository はタスクの永続化を抽象化するインターフェース。
// Save は新規作成と更新の両方を処理する（upsert）。
// FindByID に userID は不要（Entry 経由でアクセス制御を行うため）。
type TaskRepository interface {
	Save(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, id entity.TaskID) (*entity.Task, error)
	ListByEntryID(ctx context.Context, entryID entity.EntryID) ([]*entity.Task, error)
	ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error)
	Delete(ctx context.Context, id entity.TaskID) error
}
