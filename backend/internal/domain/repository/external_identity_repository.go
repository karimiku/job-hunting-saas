package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// ExternalIdentityRepository は外部認証連携情報の永続化を抽象化するインターフェース。
type ExternalIdentityRepository interface {
	Save(ctx context.Context, identity *entity.ExternalIdentity) error
	FindByProviderAndSubject(ctx context.Context, provider value.AuthProvider, subject string) (*entity.ExternalIdentity, error)
}
