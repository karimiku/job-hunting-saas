package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// PasswordCredentialRepository はパスワード認証情報の永続化を抽象化するインターフェース。
type PasswordCredentialRepository interface {
	Save(ctx context.Context, credential *entity.PasswordCredential) error
	FindByUserID(ctx context.Context, userID entity.UserID) (*entity.PasswordCredential, error)
}
