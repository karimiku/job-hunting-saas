package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// StageHistoryRepository は選考履歴の永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// StageHistory はイミュータブルなエンティティのため Save は提供せず、Create のみとする。
// 監査ログとしての性質上、Delete・FindByID も提供しない。
type StageHistoryRepository interface {
	Create(ctx context.Context, stageHistory *entity.StageHistory) error
	ListByEntryID(ctx context.Context, entryID entity.EntryID) ([]*entity.StageHistory, error)
}
