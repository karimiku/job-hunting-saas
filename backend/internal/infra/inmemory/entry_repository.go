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

// NewEntryRepository は EntryRepository を新規生成する。
func NewEntryRepository() *EntryRepository {
	return &EntryRepository{
		entriesByID: make(map[entity.EntryID]*entity.Entry),
	}
}

// Save は Entry を upsert する。同じ ID があれば更新、なければ作成。
func (r *EntryRepository) Save(_ context.Context, entry *entity.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entriesByID[entry.ID()] = entry
	return nil
}

// FindByID は userID 所有の Entry を ID から取得する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *EntryRepository) FindByID(_ context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, exists := r.entriesByID[id]
	if !exists || stored.UserID() != userID {
		return nil, repository.ErrNotFound
	}
	return stored, nil
}

// ListByUserID は userID に紐づく Entry を filter で絞り込んで返す。filter の各項目は nil なら絞り込まない。
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

// Delete は userID 所有の Entry を ID から削除する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
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
