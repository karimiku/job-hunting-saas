package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// EntryRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// 本番ではPostgreSQL実装に差し替える。
type EntryRepository struct {
	mu          sync.RWMutex
	entriesByID map[entity.EntryID]*entity.Entry
}

func NewEntryRepository() *EntryRepository {
	return &EntryRepository{
		entriesByID: make(map[entity.EntryID]*entity.Entry),
	}
}

func (r *EntryRepository) Save(_ context.Context, entry *entity.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entriesByID[entry.ID()] = entry
	return nil
}

func (r *EntryRepository) FindByID(_ context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, exists := r.entriesByID[id]
	if !exists || stored.UserID() != userID {
		return nil, repository.ErrNotFound
	}
	return stored, nil
}

func (r *EntryRepository) ListByUserID(_ context.Context, userID entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Entry
	for _, entry := range r.entriesByID {
		if entry.UserID() != userID {
			continue
		}
		if filter.Status != nil && !entry.Status().Equals(*filter.Status) {
			continue
		}
		if filter.StageKind != nil && !entry.Stage().Kind().Equals(*filter.StageKind) {
			continue
		}
		if filter.Source != nil && !entry.Source().Equals(*filter.Source) {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

func (r *EntryRepository) Delete(_ context.Context, userID entity.UserID, id entity.EntryID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, exists := r.entriesByID[id]
	if !exists || stored.UserID() != userID {
		return repository.ErrNotFound
	}
	delete(r.entriesByID, id)
	return nil
}
