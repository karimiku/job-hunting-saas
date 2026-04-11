package repository

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// TaskRepository はタスクの永続化を抽象化するインターフェース。
// Save は新規作成と更新の両方を処理する（upsert）。
// Task は UserID を持たないため、FindByID / ListByEntryID / Delete は
// 実装側で Entry を JOIN し userID によるアクセス制御を行う。
type TaskRepository interface {
	Save(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, userID entity.UserID, id entity.TaskID) (*entity.Task, error)
	ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error)
	ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error)
	Delete(ctx context.Context, userID entity.UserID, id entity.TaskID) error
}
