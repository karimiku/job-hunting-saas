package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// InboxClipRepository は Chrome 拡張クリップの永続化を抽象化するインターフェース。
// クリップは不変エンティティ — Update は提供せず、保存・取得・削除のみ。
type InboxClipRepository interface {
	Create(ctx context.Context, clip *entity.InboxClip) error
	FindByID(ctx context.Context, userID entity.UserID, id entity.InboxClipID) (*entity.InboxClip, error)
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.InboxClip, error)
	Delete(ctx context.Context, userID entity.UserID, id entity.InboxClipID) error
}
