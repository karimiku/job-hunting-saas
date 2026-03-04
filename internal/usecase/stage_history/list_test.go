package stagehistory

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestHistory(t *testing.T) *entity.StageHistory {
	t.Helper()
	stage := value.MustNewStage(value.StageKindInterview(), "一次面接")
	return entity.NewStageHistory(entity.NewEntryID(), stage, "メモ")
}

func TestList_Multiple(t *testing.T) {
	h1 := newTestHistory(t)
	h2 := newTestHistory(t)
	entryID := entity.NewEntryID()
	historyRepo := &mockHistoryRepo{
		listByEntryFn: func(_ context.Context, _ entity.EntryID) ([]*entity.StageHistory, error) {
			return []*entity.StageHistory{h1, h2}, nil
		},
	}

	uc := NewList(historyRepo, entryFound())
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entryID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.StageHistories) != 2 {
		t.Fatalf("len(StageHistories) = %d, want 2", len(out.StageHistories))
	}
}

func TestList_Empty(t *testing.T) {
	historyRepo := &mockHistoryRepo{}

	uc := NewList(historyRepo, entryFound())
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.StageHistories) != 0 {
		t.Fatalf("len(StageHistories) = %d, want 0", len(out.StageHistories))
	}
}

func TestList_EntryNotFound(t *testing.T) {
	uc := NewList(&mockHistoryRepo{}, &mockEntryRepo{})

	_, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestList_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	historyRepo := &mockHistoryRepo{
		listByEntryFn: func(_ context.Context, _ entity.EntryID) ([]*entity.StageHistory, error) {
			return nil, dbErr
		},
	}

	uc := NewList(historyRepo, entryFound())
	_, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
