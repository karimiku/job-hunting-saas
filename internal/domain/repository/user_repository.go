package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// UserRepository はユーザーの永続化を抽象化するインターフェース。
// Save は新規作成と更新の両方を処理する（upsert）。
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id entity.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email value.Email) (*entity.User, error)
	Delete(ctx context.Context, id entity.UserID) error
}
