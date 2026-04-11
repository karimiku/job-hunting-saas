package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// EntryFilter はエントリー一覧取得時のフィルタ条件。
// nil のフィールドはフィルタを適用しない。
type EntryFilter struct {
	Status    *value.EntryStatus
	StageKind *value.StageKind
	Source    *value.Source
}

// EntryRepository はエントリー（応募）の永続化を抽象化するインターフェース。
// Save は新規作成と更新の両方を処理する（upsert）。
// Save 以外のメソッドは、他ユーザーのデータを操作できないよう userID でスコープする。
type EntryRepository interface {
	Save(ctx context.Context, entry *entity.Entry) error
	FindByID(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error)
	ListByUserID(ctx context.Context, userID entity.UserID, filter EntryFilter) ([]*entity.Entry, error)
	Delete(ctx context.Context, userID entity.UserID, id entity.EntryID) error
}
