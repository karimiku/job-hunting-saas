package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// StageHistoryRepository は選考履歴の永続化を抽象化するインターフェース。
// StageHistory は不変エンティティのため Save は提供せず、Create のみとする。
// 監査ログとしての性質上 Delete・FindByID も不要。
// ListByEntryID は Entry 経由で userID の所有権を検証する。
type StageHistoryRepository interface {
	Create(ctx context.Context, history *entity.StageHistory) error
	ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.StageHistory, error)
}
