//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres"
)

func insertCommittedTestUser(t *testing.T, pool *pgxpool.Pool) entity.UserID {
	t.Helper()
	ctx := context.Background()
	userID := entity.NewUserID()

	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, email, name, created_at, updated_at)
		 VALUES ($1, $2, $3, now(), now())`,
		uuid.UUID(userID),
		uuid.New().String()+"@test.example.com",
		"テストユーザー",
	)
	if err != nil {
		t.Fatalf("failed to insert committed test user: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, uuid.UUID(userID))
	})

	return userID
}

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
