package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// ESMemoRepository は ES / 自己PR / 面接ネタ用メモの永続化を抽象化する。
type ESMemoRepository interface {
	Save(ctx context.Context, memo *entity.ESMemo) error
	ListByUserID(ctx context.Context, userID entity.UserID, limit int32) ([]*entity.ESMemo, error)
}
