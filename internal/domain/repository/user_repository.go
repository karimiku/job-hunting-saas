package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// UserRepository はユーザーの永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// Save は新規作成と更新の両方を処理する（upsert）。
// Userはシステムのルートエンティティのため、他ユーザーのスコープ制約は不要。
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, userID entity.UserID) (*entity.User, error)
	FindByEmail(ctx context.Context, email value.Email) (*entity.User, error)
	Delete(ctx context.Context, userID entity.UserID) error
}
