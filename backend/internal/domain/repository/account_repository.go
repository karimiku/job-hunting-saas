package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// AccountRepository はユーザーアカウント単位の削除を抽象化する。
type AccountRepository interface {
	DeleteAccount(ctx context.Context, userID entity.UserID) error
}
