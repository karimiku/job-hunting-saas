package inmemory

import (
	"context"
	"sort"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// CompanyAliasRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// 本番ではPostgreSQL実装に差し替える。
// 本番実装と同じ userID スコープを守り、テストが認可前提をすり抜けないようにする。
type CompanyAliasRepository struct {
	mu          sync.RWMutex
	aliasesByID map[entity.CompanyAliasID]*entity.CompanyAlias
}

// NewCompanyAliasRepository は CompanyAliasRepository を新規生成する。
func NewCompanyAliasRepository() *CompanyAliasRepository {
	return &CompanyAliasRepository{
		aliasesByID: make(map[entity.CompanyAliasID]*entity.CompanyAlias),
	}
}

// Create は CompanyAlias を保存する。同じ ID は上書きしない（不変エンティティ）。
func (r *CompanyAliasRepository) Create(_ context.Context, alias *entity.CompanyAlias) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.aliasesByID[alias.ID()]; exists {
		return repository.ErrAlreadyExists
	}
	for _, stored := range r.aliasesByID {
		if stored.UserID() == alias.UserID() &&
			stored.CompanyID() == alias.CompanyID() &&
			stored.Alias().Equals(alias.Alias()) {
			return repository.ErrAlreadyExists
		}
	}
	r.aliasesByID[alias.ID()] = alias
	return nil
}

// FindByID は userID 所有の CompanyAlias を ID から取得する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyAliasRepository) FindByID(_ context.Context, userID entity.UserID, id entity.CompanyAliasID) (*entity.CompanyAlias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, exists := r.aliasesByID[id]
	if !exists || stored.UserID() != userID {
		return nil, repository.ErrNotFound
	}
	return stored, nil
}

// ListByCompanyID は userID 所有かつ companyID に紐づく CompanyAlias を作成日時の新しい順で返す。
func (r *CompanyAliasRepository) ListByCompanyID(_ context.Context, userID entity.UserID, companyID entity.CompanyID) ([]*entity.CompanyAlias, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	owned := make([]*entity.CompanyAlias, 0)
	for _, a := range r.aliasesByID {
		if a.UserID() == userID && a.CompanyID() == companyID {
			owned = append(owned, a)
		}
	}
	sort.Slice(owned, func(i, j int) bool {
		return owned[i].CreatedAt().After(owned[j].CreatedAt())
	})
	return owned, nil
}

// Delete は userID 所有の CompanyAlias を ID から削除する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyAliasRepository) Delete(_ context.Context, userID entity.UserID, id entity.CompanyAliasID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.aliasesByID[id]
	if !exists || stored.UserID() != userID {
		return repository.ErrNotFound
	}
	delete(r.aliasesByID, id)
	return nil
}
