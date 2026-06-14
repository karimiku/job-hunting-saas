//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func TestEntryWithCompanyRepository_SaveEntryWithCompany_Success(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()
	userID := insertCommittedTestUser(t, pool)
	repo := postgres.NewEntryWithCompanyRepository(pool)

	company := entity.NewCompany(userID, newTestCompanyName(t, "同時作成企業"))
	entry := entity.NewEntry(userID, company.ID(), newTestRoute(t, "本選考"), newTestSource(t, "リクナビ"))

	if err := repo.SaveEntryWithCompany(ctx, company, entry); err != nil {
		t.Fatalf("SaveEntryWithCompany failed: %v", err)
	}

	if _, err := postgres.NewCompanyRepository(pool).FindByID(ctx, userID, company.ID()); err != nil {
		t.Fatalf("company should exist: %v", err)
	}
	if _, err := postgres.NewEntryRepository(pool).FindByID(ctx, userID, entry.ID()); err != nil {
		t.Fatalf("entry should exist: %v", err)
	}
}

func TestEntryWithCompanyRepository_SaveEntryWithCompany_RollsBackCompanyWhenEntryFails(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()
	userID := insertCommittedTestUser(t, pool)
	repo := postgres.NewEntryWithCompanyRepository(pool)

	company := entity.NewCompany(userID, newTestCompanyName(t, "ロールバック企業"))
	entry := entity.NewEntry(userID, entity.NewCompanyID(), newTestRoute(t, "本選考"), newTestSource(t, "リクナビ"))

	if err := repo.SaveEntryWithCompany(ctx, company, entry); err == nil {
		t.Fatal("SaveEntryWithCompany should fail when entry references a missing company")
	}

	_, err := postgres.NewCompanyRepository(pool).FindByID(ctx, userID, company.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("company should have been rolled back, err = %v", err)
	}
}
