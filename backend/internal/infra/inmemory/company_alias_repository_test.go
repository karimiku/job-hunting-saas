package inmemory_test

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
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
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	userID := entity.NewUserID()
	companyID := entity.NewCompanyID()
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
	if got.Alias().String() != "トヨタ" {
		t.Errorf("Alias = %q, want %q", got.Alias().String(), "トヨタ")
	}
}

func TestCompanyAliasRepository_FindByID_OtherUser(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	owner := entity.NewUserID()
	other := entity.NewUserID()
	companyID := entity.NewCompanyID()
	alias := entity.NewCompanyAlias(owner, companyID, newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err := repo.FindByID(ctx, other, alias.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestCompanyAliasRepository_Create_DuplicateAliasForSameCompany(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	userID := entity.NewUserID()
	companyID := entity.NewCompanyID()
	first := entity.NewCompanyAlias(userID, companyID, newTestAlias(t, "トヨタ"))
	second := entity.NewCompanyAlias(userID, companyID, newTestAlias(t, "トヨタ"))

	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("first Create failed: %v", err)
	}
	if err := repo.Create(ctx, second); !errors.Is(err, repository.ErrAlreadyExists) {
		t.Fatalf("second Create err = %v, want ErrAlreadyExists", err)
	}
}

func TestCompanyAliasRepository_ListByCompanyID(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	userID := entity.NewUserID()
	companyA := entity.NewCompanyID()
	companyB := entity.NewCompanyID()

	a1 := entity.NewCompanyAlias(userID, companyA, newTestAlias(t, "トヨタ"))
	a2 := entity.NewCompanyAlias(userID, companyA, newTestAlias(t, "TOYOTA"))
	b1 := entity.NewCompanyAlias(userID, companyB, newTestAlias(t, "ホンダ"))
	for _, a := range []*entity.CompanyAlias{a1, a2, b1} {
		if err := repo.Create(ctx, a); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
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

func TestCompanyAliasRepository_ListByCompanyID_ExcludesOtherUser(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	owner := entity.NewUserID()
	other := entity.NewUserID()
	companyID := entity.NewCompanyID()

	if err := repo.Create(ctx, entity.NewCompanyAlias(other, companyID, newTestAlias(t, "他人の別名"))); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.ListByCompanyID(ctx, owner, companyID)
	if err != nil {
		t.Fatalf("ListByCompanyID failed: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0 (must not leak other user data)", len(got))
	}
}

func TestCompanyAliasRepository_Delete(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	userID := entity.NewUserID()
	alias := entity.NewCompanyAlias(userID, entity.NewCompanyID(), newTestAlias(t, "トヨタ"))
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
	ctx := context.Background()
	repo := inmemory.NewCompanyAliasRepository()

	owner := entity.NewUserID()
	other := entity.NewUserID()
	alias := entity.NewCompanyAlias(owner, entity.NewCompanyID(), newTestAlias(t, "トヨタ"))
	if err := repo.Create(ctx, alias); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(ctx, other, alias.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}
