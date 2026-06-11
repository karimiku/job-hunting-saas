package inmemory

import (
	"context"
	"sort"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// ESMemoRepository はメモリ上に保存するESメモ用リポジトリ。テスト・開発用。
type ESMemoRepository struct {
	mu        sync.RWMutex
	memosByID map[entity.ESMemoID]*entity.ESMemo
}

// NewESMemoRepository は ESMemoRepository を新規生成する。
func NewESMemoRepository() *ESMemoRepository {
	return &ESMemoRepository{
		memosByID: make(map[entity.ESMemoID]*entity.ESMemo),
	}
}

// Save はESメモを保存する。
func (r *ESMemoRepository) Save(_ context.Context, memo *entity.ESMemo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.memosByID[memo.ID()] = memo
	return nil
}

// ListByUserID はユーザーのESメモを新しい順に取得する。
func (r *ESMemoRepository) ListByUserID(_ context.Context, userID entity.UserID, limit int32) ([]*entity.ESMemo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	owned := make([]*entity.ESMemo, 0)
	for _, memo := range r.memosByID {
		if memo.UserID() == userID {
			owned = append(owned, memo)
		}
	}
	sort.Slice(owned, func(i, j int) bool {
		return owned[i].CreatedAt().After(owned[j].CreatedAt())
	})
	if limit > 0 && int(limit) < len(owned) {
		owned = owned[:limit]
	}
	return owned, nil
}
