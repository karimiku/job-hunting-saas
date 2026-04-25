package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// UserRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// 本番では PostgreSQL 実装に差し替える。
type UserRepository struct {
	mu        sync.RWMutex
	usersByID map[entity.UserID]*entity.User
}

// NewUserRepository は UserRepository を新規生成する。
func NewUserRepository() *UserRepository {
	return &UserRepository{
		usersByID: make(map[entity.UserID]*entity.User),
	}
}

// Save は User を upsert する。同じ ID があれば更新、なければ作成。
func (r *UserRepository) Save(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.usersByID[user.ID()] = user
	return nil
}

// FindByID は ID から User を取得する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) FindByID(_ context.Context, id entity.UserID) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByID[id]
	if !exists {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

// FindByEmail は email から User を取得する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) FindByEmail(_ context.Context, email value.Email) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.usersByID {
		if user.Email().String() == email.String() {
			return user, nil
		}
	}
	return nil, repository.ErrNotFound
}

// Delete は ID から User を削除する。存在しない場合は repository.ErrNotFound を返す。
func (r *UserRepository) Delete(_ context.Context, id entity.UserID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.usersByID[id]; !exists {
		return repository.ErrNotFound
	}
	delete(r.usersByID, id)
	return nil
}
