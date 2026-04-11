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

func TestExternalIdentityRepository_Save_and_Find(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewExternalIdentityRepository(tx)

	identity := entity.NewExternalIdentity(userID, value.AuthProviderGoogle(), "google-sub-12345")

	if err := repo.Save(ctx, identity); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByProviderAndSubject(ctx, value.AuthProviderGoogle(), "google-sub-12345")
	if err != nil {
		t.Fatalf("FindByProviderAndSubject failed: %v", err)
	}

	if got.ID() != identity.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), identity.ID())
	}
	if got.UserID() != userID {
		t.Errorf("UserID = %v, want %v", got.UserID(), userID)
	}
	if got.Provider().String() != "google" {
		t.Errorf("Provider = %v, want google", got.Provider())
	}
	if got.Subject() != "google-sub-12345" {
		t.Errorf("Subject = %v, want google-sub-12345", got.Subject())
	}
}

func TestExternalIdentityRepository_Save_Duplicate(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewExternalIdentityRepository(tx)

	id1 := entity.NewExternalIdentity(userID, value.AuthProviderGoogle(), "dup-subject")
	if err := repo.Save(ctx, id1); err != nil {
		t.Fatalf("Save first failed: %v", err)
	}

	id2 := entity.NewExternalIdentity(userID, value.AuthProviderGoogle(), "dup-subject")
	err := repo.Save(ctx, id2)
	if !errors.Is(err, repository.ErrAlreadyExists) {
		t.Errorf("Save duplicate: got %v, want ErrAlreadyExists", err)
	}
}

func TestExternalIdentityRepository_FindByProviderAndSubject_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewExternalIdentityRepository(tx)

	_, err := repo.FindByProviderAndSubject(ctx, value.AuthProviderGoogle(), "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("FindByProviderAndSubject not found: got %v, want ErrNotFound", err)
	}
}
