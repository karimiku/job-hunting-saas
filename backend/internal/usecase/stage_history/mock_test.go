package stagehistory

import (
	"context"
	"testing"

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

// expectFindByID は FindByID が指定の (userID, entryID) で呼ばれることを検証する mockEntryRepo を返す。
// stage_history ユースケースは Entry の所有権チェック目的で FindByID を呼ぶため、
// その契約をテストで担保する。
func expectFindByID(t *testing.T, wantUserID entity.UserID, wantEntryID entity.EntryID) *mockEntryRepo {
	t.Helper()
	return &mockEntryRepo{
		findByIDFn: func(_ context.Context, gotUserID entity.UserID, gotEntryID entity.EntryID) (*entity.Entry, error) {
			if gotUserID != wantUserID {
				t.Errorf("FindByID userID = %v, want %v", gotUserID, wantUserID)
			}
			if gotEntryID != wantEntryID {
				t.Errorf("FindByID entryID = %v, want %v", gotEntryID, wantEntryID)
			}
			return &entity.Entry{}, nil
		},
	}
}

// failOnCallHistoryRepo は Create / ListByEntryID が呼ばれたら即テスト失敗する historyRepo を返す。
// EntryNotFound のときに historyRepo が呼ばれないことを担保するために使う（責務境界の保証）。
func failOnCallHistoryRepo(t *testing.T) *mockHistoryRepo {
	t.Helper()
	return &mockHistoryRepo{
		createFn: func(_ context.Context, _ *entity.StageHistory) error {
			t.Error("historyRepo.Create should not be called when entry lookup fails")
			return nil
		},
		listFn: func(_ context.Context, _ entity.EntryID) ([]*entity.StageHistory, error) {
			t.Error("historyRepo.ListByEntryID should not be called when entry lookup fails")
			return nil, nil
		},
	}
}
