package stagehistory

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestList_Success(t *testing.T) {
	entryID := entity.NewEntryID()
	histories := []*entity.StageHistory{
		entity.NewStageHistory(entryID, value.MustNewStage(value.StageKindDocument(), "ES提出"), ""),
		entity.NewStageHistory(entryID, value.MustNewStage(value.StageKindInterview(), "一次面接"), "オンライン"),
	}

	historyRepo := &mockHistoryRepo{
		listFn: func(_ context.Context, _ entity.EntryID) ([]*entity.StageHistory, error) {
			return histories, nil
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
		t.Errorf("len = %d, want 2", len(out.StageHistories))
	}
}

func TestList_Empty(t *testing.T) {
	uc := NewList(&mockHistoryRepo{}, entryFound())
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.StageHistories) != 0 {
		t.Errorf("len = %d, want 0", len(out.StageHistories))
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

func TestList_RepoError(t *testing.T) {
	listErr := errors.New("db read failed")
	historyRepo := &mockHistoryRepo{
		listFn: func(_ context.Context, _ entity.EntryID) ([]*entity.StageHistory, error) {
			return nil, listErr
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
	if !errors.Is(err, listErr) {
		t.Errorf("error = %v, want listErr", err)
	}
}
