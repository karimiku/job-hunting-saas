//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func newTestRoute(t *testing.T, raw string) value.Route {
	t.Helper()
	r, err := value.NewRoute(raw)
	if err != nil {
		t.Fatalf("NewRoute(%q) failed: %v", raw, err)
	}
	return r
}

func newTestSource(t *testing.T, raw string) value.Source {
	t.Helper()
	s, err := value.NewSource(raw)
	if err != nil {
		t.Fatalf("NewSource(%q) failed: %v", raw, err)
	}
	return s
}

func newTestEntryStatus(t *testing.T, raw string) value.EntryStatus {
	t.Helper()
	s, err := value.NewEntryStatus(raw)
	if err != nil {
		t.Fatalf("NewEntryStatus(%q) failed: %v", raw, err)
	}
	return s
}

func newTestStageKind(t *testing.T, raw string) value.StageKind {
	t.Helper()
	k, err := value.NewStageKind(raw)
	if err != nil {
		t.Fatalf("NewStageKind(%q) failed: %v", raw, err)
	}
	return k
}

func newTestStage(t *testing.T, kindRaw, label string) value.Stage {
	t.Helper()
	kind := newTestStageKind(t, kindRaw)
	s, err := value.NewStage(kind, label)
	if err != nil {
		t.Fatalf("NewStage(%q, %q) failed: %v", kindRaw, label, err)
	}
	return s
}

func TestEntryRepository_Save_Insert(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	entry := entity.NewEntry(userID, companyID, newTestRoute(t, "直接応募"), newTestSource(t, "マイナビ"))

	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, entry.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.ID() != entry.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), entry.ID())
	}
	if got.Route().String() != "直接応募" {
		t.Errorf("Route = %q, want %q", got.Route().String(), "直接応募")
	}
	if got.Source().String() != "マイナビ" {
		t.Errorf("Source = %q, want %q", got.Source().String(), "マイナビ")
	}
	if got.Status().String() != "in_progress" {
		t.Errorf("Status = %q, want %q", got.Status().String(), "in_progress")
	}
}

func TestEntryRepository_Save_Update(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	entry := entity.NewEntry(userID, companyID, newTestRoute(t, "直接応募"), newTestSource(t, "マイナビ"))
	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save (insert) failed: %v", err)
	}

	entry.UpdateStatus(newTestEntryStatus(t, "offered"))
	entry.UpdateStage(newTestStage(t, "offer", "内定"))
	entry.UpdateMemo("面接合格")
	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save (update) failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, entry.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.Status().String() != "offered" {
		t.Errorf("Status = %q, want %q", got.Status().String(), "offered")
	}
	if got.Stage().Kind().String() != "offer" {
		t.Errorf("Stage.Kind = %q, want %q", got.Stage().Kind().String(), "offer")
	}
	if got.Stage().Label() != "内定" {
		t.Errorf("Stage.Label = %q, want %q", got.Stage().Label(), "内定")
	}
	if got.Memo() != "面接合格" {
		t.Errorf("Memo = %q, want %q", got.Memo(), "面接合格")
	}
}

func TestEntryRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewEntryRepository(tx)

	_, err := repo.FindByID(ctx, userID, entity.NewEntryID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestEntryRepository_FindByID_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	entry := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "リクナビ"))
	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err := repo.FindByID(ctx, otherUserID, entry.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestEntryRepository_ListByUserID_NoFilter(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	e1 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	e2 := entity.NewEntry(userID, companyID, newTestRoute(t, "エージェント"), newTestSource(t, "リクナビ"))
	for _, e := range []*entity.Entry{e1, e2} {
		if err := repo.Save(ctx, e); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	list, err := repo.ListByUserID(ctx, userID, repository.EntryFilter{})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestEntryRepository_ListByUserID_FilterByStatus(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	e1 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	e2 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "リクナビ"))
	e2.UpdateStatus(newTestEntryStatus(t, "offered"))
	for _, e := range []*entity.Entry{e1, e2} {
		if err := repo.Save(ctx, e); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	status := newTestEntryStatus(t, "offered")
	list, err := repo.ListByUserID(ctx, userID, repository.EntryFilter{Status: &status})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != e2.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), e2.ID())
	}
}

func TestEntryRepository_ListByUserID_FilterByStageKind(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	e1 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	e2 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "リクナビ"))
	e2.UpdateStage(newTestStage(t, "interview", "一次面接"))
	for _, e := range []*entity.Entry{e1, e2} {
		if err := repo.Save(ctx, e); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	kind := newTestStageKind(t, "interview")
	list, err := repo.ListByUserID(ctx, userID, repository.EntryFilter{StageKind: &kind})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != e2.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), e2.ID())
	}
}

func TestEntryRepository_ListByUserID_FilterBySource(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	e1 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	e2 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "リクナビ"))
	for _, e := range []*entity.Entry{e1, e2} {
		if err := repo.Save(ctx, e); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	source := newTestSource(t, "リクナビ")
	list, err := repo.ListByUserID(ctx, userID, repository.EntryFilter{Source: &source})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != e2.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), e2.ID())
	}
}

func TestEntryRepository_ListByUserID_FilterCombined(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	// e1: in_progress, application, マイナビ
	e1 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	// e2: offered, interview, マイナビ
	e2 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	e2.UpdateStatus(newTestEntryStatus(t, "offered"))
	e2.UpdateStage(newTestStage(t, "interview", "一次面接"))
	// e3: offered, interview, リクナビ
	e3 := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "リクナビ"))
	e3.UpdateStatus(newTestEntryStatus(t, "offered"))
	e3.UpdateStage(newTestStage(t, "interview", "一次面接"))

	for _, e := range []*entity.Entry{e1, e2, e3} {
		if err := repo.Save(ctx, e); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	status := newTestEntryStatus(t, "offered")
	kind := newTestStageKind(t, "interview")
	source := newTestSource(t, "マイナビ")
	list, err := repo.ListByUserID(ctx, userID, repository.EntryFilter{
		Status:    &status,
		StageKind: &kind,
		Source:    &source,
	})
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != e2.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), e2.ID())
	}
}

func TestEntryRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	entry := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := repo.Delete(ctx, userID, entry.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, userID, entry.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("after delete: err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestEntryRepository_Delete_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewEntryRepository(tx)

	entry := entity.NewEntry(userID, companyID, newTestRoute(t, "直接"), newTestSource(t, "マイナビ"))
	if err := repo.Save(ctx, entry); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	err := repo.Delete(ctx, otherUserID, entry.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}
