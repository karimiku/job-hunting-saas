package stagehistory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type mockEntryRepo struct {
	findByIDFn func(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error)
}

func (m *mockEntryRepo) Save(_ context.Context, _ *entity.Entry) error {
	return nil
}

func (m *mockEntryRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockEntryRepo) ListByUserID(_ context.Context, _ entity.UserID, _ repository.EntryFilter) ([]*entity.Entry, error) {
	return []*entity.Entry{}, nil
}

func (m *mockEntryRepo) Delete(_ context.Context, _ entity.UserID, _ entity.EntryID) error {
	return nil
}

type mockHistoryRepo struct {
	createFn func(ctx context.Context, history *entity.StageHistory) error
	listFn   func(ctx context.Context, entryID entity.EntryID) ([]*entity.StageHistory, error)
}

func (m *mockHistoryRepo) Create(ctx context.Context, history *entity.StageHistory) error {
	if m.createFn != nil {
		return m.createFn(ctx, history)
	}
	return nil
}

func (m *mockHistoryRepo) ListByEntryID(ctx context.Context, entryID entity.EntryID) ([]*entity.StageHistory, error) {
	if m.listFn != nil {
		return m.listFn(ctx, entryID)
	}
	return []*entity.StageHistory{}, nil
}

func entryFound() *mockEntryRepo {
	return &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return &entity.Entry{}, nil
		},
	}
}
