//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func newTestTaskTitle(t *testing.T, raw string) value.TaskTitle {
	t.Helper()
	title, err := value.NewTaskTitle(raw)
	if err != nil {
		t.Fatalf("NewTaskTitle(%q) failed: %v", raw, err)
	}
	return title
}

func newTestTaskType(t *testing.T, raw string) value.TaskType {
	t.Helper()
	tt, err := value.NewTaskType(raw)
	if err != nil {
		t.Fatalf("NewTaskType(%q) failed: %v", raw, err)
	}
	return tt
}

func TestTaskRepository_Save_Insert(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	task := entity.NewTask(entryID, newTestTaskTitle(t, "ES提出"), newTestTaskType(t, "deadline"))

	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, task.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.ID() != task.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), task.ID())
	}
	if got.Title().String() != "ES提出" {
		t.Errorf("Title = %q, want %q", got.Title().String(), "ES提出")
	}
	if got.TaskType().String() != "deadline" {
		t.Errorf("TaskType = %q, want %q", got.TaskType().String(), "deadline")
	}
	if got.Status().String() != "todo" {
		t.Errorf("Status = %q, want %q", got.Status().String(), "todo")
	}
	if got.DueDate() != nil {
		t.Errorf("DueDate = %v, want nil", got.DueDate())
	}
}

func TestTaskRepository_Save_Update(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	task := entity.NewTask(entryID, newTestTaskTitle(t, "ES提出"), newTestTaskType(t, "deadline"))
	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save (insert) failed: %v", err)
	}

	due := time.Date(2026, 4, 15, 23, 59, 0, 0, time.UTC)
	task.SetDueDate(due)
	task.Complete()
	task.UpdateMemo("提出済み")
	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save (update) failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, task.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.Status().String() != "done" {
		t.Errorf("Status = %q, want %q", got.Status().String(), "done")
	}
	if got.DueDate() == nil {
		t.Fatal("DueDate should not be nil")
	}
	if !got.DueDate().Truncate(time.Microsecond).Equal(due.Truncate(time.Microsecond)) {
		t.Errorf("DueDate = %v, want %v", got.DueDate(), due)
	}
	if got.Memo() != "提出済み" {
		t.Errorf("Memo = %q, want %q", got.Memo(), "提出済み")
	}
}

func TestTaskRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewTaskRepository(tx)

	_, err := repo.FindByID(ctx, userID, entity.NewTaskID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestTaskRepository_FindByID_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	task := entity.NewTask(entryID, newTestTaskTitle(t, "ES提出"), newTestTaskType(t, "deadline"))
	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err := repo.FindByID(ctx, otherUserID, task.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestTaskRepository_ListByEntryID(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	t1 := entity.NewTask(entryID, newTestTaskTitle(t, "タスク1"), newTestTaskType(t, "deadline"))
	t2 := entity.NewTask(entryID, newTestTaskTitle(t, "タスク2"), newTestTaskType(t, "schedule"))
	for _, task := range []*entity.Task{t1, t2} {
		if err := repo.Save(ctx, task); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	list, err := repo.ListByEntryID(ctx, userID, entryID)
	if err != nil {
		t.Fatalf("ListByEntryID failed: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestTaskRepository_ListByUserIDWithDueBefore(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	deadline := time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC)

	// due_dateが期限前（未完了）→ 含まれる
	t1 := entity.NewTask(entryID, newTestTaskTitle(t, "期限前未完了"), newTestTaskType(t, "deadline"))
	t1.SetDueDate(time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC))
	// due_dateが期限後
	t2 := entity.NewTask(entryID, newTestTaskTitle(t, "期限後"), newTestTaskType(t, "deadline"))
	t2.SetDueDate(time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC))
	// due_dateなし
	t3 := entity.NewTask(entryID, newTestTaskTitle(t, "期限なし"), newTestTaskType(t, "schedule"))
	// due_dateが期限前だが完了済み → 含まれない
	t4 := entity.NewTask(entryID, newTestTaskTitle(t, "期限前完了済み"), newTestTaskType(t, "deadline"))
	t4.SetDueDate(time.Date(2026, 4, 3, 0, 0, 0, 0, time.UTC))
	t4.Complete()

	for _, task := range []*entity.Task{t1, t2, t3, t4} {
		if err := repo.Save(ctx, task); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	list, err := repo.ListByUserIDWithDueBefore(ctx, userID, deadline)
	if err != nil {
		t.Fatalf("ListByUserIDWithDueBefore failed: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != t1.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), t1.ID())
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	task := entity.NewTask(entryID, newTestTaskTitle(t, "削除対象"), newTestTaskType(t, "deadline"))
	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := repo.Delete(ctx, userID, task.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, userID, task.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("after delete: err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestTaskRepository_Delete_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	entryID := insertTestEntry(t, tx, userID, companyID)
	repo := postgres.NewTaskRepository(tx)

	task := entity.NewTask(entryID, newTestTaskTitle(t, "削除対象"), newTestTaskType(t, "deadline"))
	if err := repo.Save(ctx, task); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	err := repo.Delete(ctx, otherUserID, task.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}
