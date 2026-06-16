package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// SelectionFlowRepository はEntryごとの可変選考フローの永続化を抽象化する。
type SelectionFlowRepository interface {
	Upsert(ctx context.Context, flow *entity.SelectionFlow) (*entity.SelectionFlow, error)
	FindByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) (*entity.SelectionFlow, error)
}
