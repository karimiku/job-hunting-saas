package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// InboxClipRepository は Chrome 拡張クリップの永続化を抽象化するインターフェース。
// クリップは不変エンティティ — Update は提供せず、保存・取得・削除のみ。
type InboxClipRepository interface {
	Create(ctx context.Context, clip *entity.InboxClip) error
	FindByID(ctx context.Context, userID entity.UserID, id entity.InboxClipID) (*entity.InboxClip, error)
	// FindByUserIDAndURL は同一ユーザーが同じ URL で保存済みのクリップを返す。
	// 存在しなければ ErrNotFound。重複クリップの抑止 (#98) に使用する。
	FindByUserIDAndURL(ctx context.Context, userID entity.UserID, url value.URL) (*entity.InboxClip, error)
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.InboxClip, error)
	Delete(ctx context.Context, userID entity.UserID, id entity.InboxClipID) error
}
