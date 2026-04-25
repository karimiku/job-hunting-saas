package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// ExternalIdentityRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// 本番では PostgreSQL 実装に差し替える。
type ExternalIdentityRepository struct {
	mu   sync.RWMutex
	byID map[entity.ExternalIdentityID]*entity.ExternalIdentity
}

// NewExternalIdentityRepository は ExternalIdentityRepository を新規生成する。
func NewExternalIdentityRepository() *ExternalIdentityRepository {
	return &ExternalIdentityRepository{
		byID: make(map[entity.ExternalIdentityID]*entity.ExternalIdentity),
	}
}

// Save は ExternalIdentity を upsert する。
func (r *ExternalIdentityRepository) Save(_ context.Context, identity *entity.ExternalIdentity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[identity.ID()] = identity
	return nil
}

// FindByProviderAndSubject は (provider, subject) のペアから ExternalIdentity を取得する。
// 存在しない場合は repository.ErrNotFound を返す。
func (r *ExternalIdentityRepository) FindByProviderAndSubject(_ context.Context, provider value.AuthProvider, subject string) (*entity.ExternalIdentity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, identity := range r.byID {
		if identity.Provider().Equals(provider) && identity.Subject() == subject {
			return identity, nil
		}
	}
	return nil, repository.ErrNotFound
}
