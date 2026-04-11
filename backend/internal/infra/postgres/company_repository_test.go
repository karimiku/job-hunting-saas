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

func newTestCompanyName(t *testing.T, raw string) value.CompanyName {
	t.Helper()
	name, err := value.NewCompanyName(raw)
	if err != nil {
		t.Fatalf("NewCompanyName(%q) failed: %v", raw, err)
	}
	return name
}

func TestCompanyRepository_Save_Insert(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	company := entity.NewCompany(userID, newTestCompanyName(t, "トヨタ自動車"))

	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, company.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.ID() != company.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), company.ID())
	}
	if got.UserID() != userID {
		t.Errorf("UserID = %v, want %v", got.UserID(), userID)
	}
	if got.Name().String() != "トヨタ自動車" {
		t.Errorf("Name = %q, want %q", got.Name().String(), "トヨタ自動車")
	}
	if got.Memo() != "" {
		t.Errorf("Memo = %q, want %q", got.Memo(), "")
	}
}

func TestCompanyRepository_Save_Update(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	company := entity.NewCompany(userID, newTestCompanyName(t, "旧社名"))
	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save (insert) failed: %v", err)
	}

	company.Rename(newTestCompanyName(t, "新社名"))
	company.UpdateMemo("メモ更新")
	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save (update) failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, company.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.Name().String() != "新社名" {
		t.Errorf("Name = %q, want %q", got.Name().String(), "新社名")
	}
	if got.Memo() != "メモ更新" {
		t.Errorf("Memo = %q, want %q", got.Memo(), "メモ更新")
	}
}

func TestCompanyRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	_, err := repo.FindByID(ctx, userID, entity.NewCompanyID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestCompanyRepository_FindByID_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	company := entity.NewCompany(userID, newTestCompanyName(t, "自社"))
	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err := repo.FindByID(ctx, otherUserID, company.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestCompanyRepository_ListByUserID(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	c1 := entity.NewCompany(userID, newTestCompanyName(t, "企業A"))
	c2 := entity.NewCompany(userID, newTestCompanyName(t, "企業B"))
	c3 := entity.NewCompany(otherUserID, newTestCompanyName(t, "他ユーザー企業"))

	for _, c := range []*entity.Company{c1, c2, c3} {
		if err := repo.Save(ctx, c); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	list, err := repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
}

func TestCompanyRepository_ListByUserID_Empty(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	list, err := repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}

	if list == nil {
		t.Error("list should not be nil")
	}
	if len(list) != 0 {
		t.Errorf("len = %d, want 0", len(list))
	}
}

func TestCompanyRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	company := entity.NewCompany(userID, newTestCompanyName(t, "削除対象"))
	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := repo.Delete(ctx, userID, company.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, userID, company.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("after delete: err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestCompanyRepository_Delete_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewCompanyRepository(tx)

	company := entity.NewCompany(userID, newTestCompanyName(t, "削除対象"))
	if err := repo.Save(ctx, company); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	err := repo.Delete(ctx, otherUserID, company.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}
