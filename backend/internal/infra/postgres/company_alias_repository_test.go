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

func newTestAlias(t *testing.T, raw string) value.Alias {
	t.Helper()
	a, err := value.NewAlias(raw)
	if err != nil {
		t.Fatalf("NewAlias(%q) failed: %v", raw, err)
	}
	return a
}

func TestCompanyAliasRepository_Create_FindByID(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewCompanyAliasRepository(tx)

	alias := entity.NewCompanyAlias(userID, companyID, newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, alias.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if got.ID() != alias.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), alias.ID())
	}
	if got.CompanyID() != companyID {
		t.Errorf("CompanyID = %v, want %v", got.CompanyID(), companyID)
	}
	if got.Alias().String() != "トヨタ" {
		t.Errorf("Alias = %q, want %q", got.Alias().String(), "トヨタ")
	}
}

func TestCompanyAliasRepository_FindByID_OtherUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	owner := insertTestUser(t, tx)
	other := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, owner)
	repo := postgres.NewCompanyAliasRepository(tx)

	alias := entity.NewCompanyAlias(owner, companyID, newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if _, err := repo.FindByID(ctx, other, alias.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestCompanyAliasRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewCompanyAliasRepository(tx)

	if _, err := repo.FindByID(ctx, userID, entity.NewCompanyAliasID()); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestCompanyAliasRepository_ListByCompanyID(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyA := insertTestCompany(t, tx, userID)
	companyB := insertTestCompany(t, tx, userID)
	repo := postgres.NewCompanyAliasRepository(tx)

	for _, raw := range []string{"トヨタ", "TOYOTA"} {
		if err := repo.Create(ctx, entity.NewCompanyAlias(userID, companyA, newTestAlias(t, raw))); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}
	if err := repo.Create(ctx, entity.NewCompanyAlias(userID, companyB, newTestAlias(t, "ホンダ"))); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.ListByCompanyID(ctx, userID, companyA)
	if err != nil {
		t.Fatalf("ListByCompanyID failed: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	for _, a := range got {
		if a.CompanyID() != companyA {
			t.Errorf("CompanyID = %v, want %v", a.CompanyID(), companyA)
		}
	}
}

func TestCompanyAliasRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, userID)
	repo := postgres.NewCompanyAliasRepository(tx)

	alias := entity.NewCompanyAlias(userID, companyID, newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(ctx, userID, alias.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := repo.FindByID(ctx, userID, alias.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("after delete err = %v, want ErrNotFound", err)
	}
}

func TestCompanyAliasRepository_Delete_OtherUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	owner := insertTestUser(t, tx)
	other := insertTestUser(t, tx)
	companyID := insertTestCompany(t, tx, owner)
	repo := postgres.NewCompanyAliasRepository(tx)

	alias := entity.NewCompanyAlias(owner, companyID, newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(ctx, other, alias.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
