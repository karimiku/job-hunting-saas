package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// TaskRepository はメモリ上にデータを保持するテスト・開発用のリポジトリ実装。
// Task は UserID を直接持たないため、EntryRepository を参照して
// Entry 経由でユーザーの所有権を検証する。
type TaskRepository struct {
	mu        sync.RWMutex
	tasksByID map[entity.TaskID]*entity.Task
	entryRepo repository.EntryRepository
}

func NewTaskRepository(entryRepo repository.EntryRepository) *TaskRepository {
	return &TaskRepository{
		tasksByID: make(map[entity.TaskID]*entity.Task),
		entryRepo: entryRepo,
	}
}

func (r *TaskRepository) Save(_ context.Context, task *entity.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasksByID[task.ID()] = task
	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.TaskID) (*entity.Task, error) {
	r.mu.RLock()
	task, exists := r.tasksByID[id]
	r.mu.RUnlock()

	if !exists {
		return nil, repository.ErrNotFound
	}

	// Entry経由でuserIDの所有権を検証
	if _, err := r.entryRepo.FindByID(ctx, userID, task.EntryID()); err != nil {
		return nil, repository.ErrNotFound
	}
	return task, nil
}

func (r *TaskRepository) ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error) {
	// Entry経由でuserIDの所有権を検証
	if _, err := r.entryRepo.FindByID(ctx, userID, entryID); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Task
	for _, task := range r.tasksByID {
		if task.EntryID() == entryID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *TaskRepository) ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Task
	for _, task := range r.tasksByID {
		if task.Status().IsDone() {
			continue
		}
		if task.DueDate() == nil || !task.DueDate().Before(deadline) {
			continue
		}
		// Entry経由でuserIDの所有権を検証
		if _, err := r.entryRepo.FindByID(ctx, userID, task.EntryID()); err != nil {
			continue
		}
		result = append(result, task)
	}
	return result, nil
}

func (r *TaskRepository) Delete(ctx context.Context, userID entity.UserID, id entity.TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasksByID[id]
	if !exists {
		return repository.ErrNotFound
	}

	// Entry経由でuserIDの所有権を検証
	if _, err := r.entryRepo.FindByID(ctx, userID, task.EntryID()); err != nil {
		return repository.ErrNotFound
	}

	delete(r.tasksByID, id)
	return nil
}
