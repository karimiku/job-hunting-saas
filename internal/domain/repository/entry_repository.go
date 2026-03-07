package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// EntryFilter はエントリー一覧取得時のフィルタ条件。
// nilのフィールドはフィルタを適用しない（全件取得と同等）。
type EntryFilter struct {
	Status    *value.EntryStatus
	StageKind *value.StageKind
	Source    *value.Source
}

// EntryRepository はエントリー（応募）の永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// Save は新規作成と更新の両方を処理する（upsert）。存在判定の責務はusecase側にある。
// Save 以外のメソッドは、他ユーザーのデータへのアクセスを防ぐため userID でスコープする。
type EntryRepository interface {
	Save(ctx context.Context, entry *entity.Entry) error
	FindByID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) (*entity.Entry, error)
	ListByUserID(ctx context.Context, userID entity.UserID, filter EntryFilter) ([]*entity.Entry, error)
	Delete(ctx context.Context, userID entity.UserID, entryID entity.EntryID) error
}
