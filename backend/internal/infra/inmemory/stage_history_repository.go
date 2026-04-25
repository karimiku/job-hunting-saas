package inmemory

import (
	"context"
	"sync"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// StageHistoryRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// StageHistory はイミュータブルのため Create のみ提供する。
type StageHistoryRepository struct {
	mu            sync.RWMutex
	historiesByID map[entity.StageHistoryID]*entity.StageHistory
}

// NewStageHistoryRepository は StageHistoryRepository を新規生成する。
func NewStageHistoryRepository() *StageHistoryRepository {
	return &StageHistoryRepository{
		historiesByID: make(map[entity.StageHistoryID]*entity.StageHistory),
	}
}

// Create は StageHistory を作成する (immutable なので Save なし)。
func (r *StageHistoryRepository) Create(_ context.Context, history *entity.StageHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.historiesByID[history.ID()] = history
	return nil
}

// ListByEntryID は entry に紐づく StageHistory を全件返す。
func (r *StageHistoryRepository) ListByEntryID(_ context.Context, entryID entity.EntryID) ([]*entity.StageHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.StageHistory
	for _, h := range r.historiesByID {
		if h.EntryID() == entryID {
			result = append(result, h)
		}
	}
	return result, nil
}
