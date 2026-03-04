package task

import (
	"context"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type mockTaskRepo struct {
	saveFn            func(ctx context.Context, task *entity.Task) error
	findByIDFn        func(ctx context.Context, userID entity.UserID, id entity.TaskID) (*entity.Task, error)
	listByEntryIDFn   func(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error)
	listByDueBeforeFn func(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error)
	deleteFn          func(ctx context.Context, userID entity.UserID, id entity.TaskID) error
}

func (m *mockTaskRepo) Save(ctx context.Context, task *entity.Task) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, task)
	}
	return nil
}

func (m *mockTaskRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.TaskID) (*entity.Task, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockTaskRepo) ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error) {
	if m.listByEntryIDFn != nil {
		return m.listByEntryIDFn(ctx, userID, entryID)
	}
	return []*entity.Task{}, nil
}

func (m *mockTaskRepo) ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error) {
	if m.listByDueBeforeFn != nil {
		return m.listByDueBeforeFn(ctx, userID, deadline)
	}
	return []*entity.Task{}, nil
}

func (m *mockTaskRepo) Delete(ctx context.Context, userID entity.UserID, id entity.TaskID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID, id)
	}
	return nil
}

type mockEntryRepo struct {
	findByIDFn func(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error)
}

func (m *mockEntryRepo) Save(ctx context.Context, entry *entity.Entry) error {
	return nil
}

func (m *mockEntryRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockEntryRepo) ListByUserID(ctx context.Context, userID entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error) {
	return []*entity.Entry{}, nil
}

func (m *mockEntryRepo) Delete(ctx context.Context, userID entity.UserID, id entity.EntryID) error {
	return nil
}
