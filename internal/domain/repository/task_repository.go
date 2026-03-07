package repository

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// TaskRepository はタスクの永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// Task は UserID を持たないため、FindByID / ListByEntryID / Delete は
// 実装側で Entry を JOIN し userID によるアクセス制御を行う。
type TaskRepository interface {
	Save(ctx context.Context, task *entity.Task) error
	FindByID(ctx context.Context, userID entity.UserID, taskID entity.TaskID) (*entity.Task, error)
	ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error)
	ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error)
	Delete(ctx context.Context, userID entity.UserID, taskID entity.TaskID) error
}
