package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// SelectionFlowRepository はメモリ上に選考フローを保存するテスト・開発用リポジトリ。
type SelectionFlowRepository struct {
	mu           sync.RWMutex
	flowsByEntry map[entity.EntryID]*entity.SelectionFlow
}

// NewSelectionFlowRepository は SelectionFlowRepository を新規生成する。
func NewSelectionFlowRepository() *SelectionFlowRepository {
	return &SelectionFlowRepository{
		flowsByEntry: make(map[entity.EntryID]*entity.SelectionFlow),
	}
}

// Upsert はEntryごとの選考フローを置換保存する。
func (r *SelectionFlowRepository) Upsert(_ context.Context, flow *entity.SelectionFlow) (*entity.SelectionFlow, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flowsByEntry[flow.EntryID()] = flow
	return flow, nil
}

// FindByEntryID はEntryに紐づく選考フローを返す。
func (r *SelectionFlowRepository) FindByEntryID(_ context.Context, _ entity.UserID, entryID entity.EntryID) (*entity.SelectionFlow, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	flow, ok := r.flowsByEntry[entryID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return flow, nil
}
