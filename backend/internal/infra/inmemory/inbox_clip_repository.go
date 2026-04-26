package inmemory

import (
	"context"
	"sort"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// InboxClipRepository はメモリ上に保存するクリップ用リポジトリ。テスト・開発用。
type InboxClipRepository struct {
	mu        sync.RWMutex
	clipsByID map[entity.InboxClipID]*entity.InboxClip
}

// NewInboxClipRepository は InboxClipRepository を新規生成する。
func NewInboxClipRepository() *InboxClipRepository {
	return &InboxClipRepository{
		clipsByID: make(map[entity.InboxClipID]*entity.InboxClip),
	}
}

// Create はクリップを保存する。同じ ID は上書きしない（不変エンティティ）。
func (r *InboxClipRepository) Create(_ context.Context, clip *entity.InboxClip) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clipsByID[clip.ID()] = clip
	return nil
}

// FindByID は userID 所有のクリップを ID から取得する。
func (r *InboxClipRepository) FindByID(_ context.Context, userID entity.UserID, id entity.InboxClipID) (*entity.InboxClip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	clip, ok := r.clipsByID[id]
	if !ok || clip.UserID() != userID {
		return nil, repository.ErrNotFound
	}
	return clip, nil
}

// ListByUserID は userID 所有のクリップを保存日時の新しい順で返す。
func (r *InboxClipRepository) ListByUserID(_ context.Context, userID entity.UserID) ([]*entity.InboxClip, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	owned := make([]*entity.InboxClip, 0)
	for _, c := range r.clipsByID {
		if c.UserID() == userID {
			owned = append(owned, c)
		}
	}
	sort.Slice(owned, func(i, j int) bool {
		return owned[i].CapturedAt().After(owned[j].CapturedAt())
	})
	return owned, nil
}

// Delete は userID 所有のクリップを削除する。所有していなければ ErrNotFound。
func (r *InboxClipRepository) Delete(_ context.Context, userID entity.UserID, id entity.InboxClipID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	clip, ok := r.clipsByID[id]
	if !ok || clip.UserID() != userID {
		return repository.ErrNotFound
	}
	delete(r.clipsByID, id)
	return nil
}
