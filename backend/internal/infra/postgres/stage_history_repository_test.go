//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func TestStageHistoryRepository_Create(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewStageHistoryRepository(tx)

	stage := newTestStage(t, "interview", "一次面接")
	history := entity.NewStageHistory(entryID, stage, "面接実施")

	if err := repo.Create(ctx, history); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	list, err := repo.ListByEntryID(ctx, entryID)
	if err != nil {
		t.Fatalf("ListByEntryID failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}

	got := list[0]
	if got.ID() != history.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), history.ID())
	}
	if got.EntryID() != entryID {
		t.Errorf("EntryID = %v, want %v", got.EntryID(), entryID)
	}
	if got.Stage().Kind().String() != "interview" {
		t.Errorf("Stage.Kind = %q, want %q", got.Stage().Kind().String(), "interview")
	}
	if got.Stage().Label() != "一次面接" {
		t.Errorf("Stage.Label = %q, want %q", got.Stage().Label(), "一次面接")
	}
	if got.Note() != "面接実施" {
		t.Errorf("Note = %q, want %q", got.Note(), "面接実施")
	}
}

func TestStageHistoryRepository_ListByEntryID_Order(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewStageHistoryRepository(tx)

	earlier := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	later := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
	h1 := entity.ReconstructStageHistory(entity.NewStageHistoryID(), entryID, newTestStage(t, "application", "応募"), "応募完了", earlier)
	h2 := entity.ReconstructStageHistory(entity.NewStageHistoryID(), entryID, newTestStage(t, "document", "書類選考"), "書類提出", later)

	for _, h := range []*entity.StageHistory{h1, h2} {
		if err := repo.Create(ctx, h); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	list, err := repo.ListByEntryID(ctx, entryID)
	if err != nil {
		t.Fatalf("ListByEntryID failed: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}

	// created_at 昇順
	if list[0].ID() != h1.ID() {
		t.Errorf("first ID = %v, want %v", list[0].ID(), h1.ID())
	}
	if list[1].ID() != h2.ID() {
		t.Errorf("second ID = %v, want %v", list[1].ID(), h2.ID())
	}
}

func TestStageHistoryRepository_ListByEntryID_Empty(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewStageHistoryRepository(tx)

	list, err := repo.ListByEntryID(ctx, entryID)
	if err != nil {
		t.Fatalf("ListByEntryID failed: %v", err)
	}

	if list == nil {
		t.Error("list should not be nil")
	}
	if len(list) != 0 {
		t.Errorf("len = %d, want 0", len(list))
	}
}
