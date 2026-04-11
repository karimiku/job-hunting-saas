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

func newTestEmail(t *testing.T, raw string) value.Email {
	t.Helper()
	email, err := value.NewEmail(raw)
	if err != nil {
		t.Fatalf("NewEmail(%q) failed: %v", raw, err)
	}
	return email
}

func newTestUserName(t *testing.T, raw string) value.UserName {
	t.Helper()
	name, err := value.NewUserName(raw)
	if err != nil {
		t.Fatalf("NewUserName(%q) failed: %v", raw, err)
	}
	return name
}

func TestUserRepository_Save_Insert(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)
	user := entity.NewUser(newTestEmail(t, "test@example.com"), newTestUserName(t, "テストユーザー"))

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.ID() != user.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), user.ID())
	}
	if got.Email().String() != "test@example.com" {
		t.Errorf("Email = %v, want test@example.com", got.Email())
	}
	if got.Name().String() != "テストユーザー" {
		t.Errorf("Name = %v, want テストユーザー", got.Name())
	}
}

func TestUserRepository_Save_Update(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)
	user := entity.NewUser(newTestEmail(t, "before@example.com"), newTestUserName(t, "変更前"))

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save (insert) failed: %v", err)
	}

	user.Rename(newTestUserName(t, "変更後"))
	user.ChangeEmail(newTestEmail(t, "after@example.com"))

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save (update) failed: %v", err)
	}

	got, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.Name().String() != "変更後" {
		t.Errorf("Name = %v, want 変更後", got.Name())
	}
	if got.Email().String() != "after@example.com" {
		t.Errorf("Email = %v, want after@example.com", got.Email())
	}
}

func TestUserRepository_Save_DuplicateEmail(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)

	user1 := entity.NewUser(newTestEmail(t, "dup@example.com"), newTestUserName(t, "ユーザー1"))
	if err := repo.Save(ctx, user1); err != nil {
		t.Fatalf("Save user1 failed: %v", err)
	}

	user2 := entity.NewUser(newTestEmail(t, "dup@example.com"), newTestUserName(t, "ユーザー2"))
	err := repo.Save(ctx, user2)
	if !errors.Is(err, repository.ErrAlreadyExists) {
		t.Errorf("Save duplicate email: got %v, want ErrAlreadyExists", err)
	}
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)

	_, err := repo.FindByID(ctx, entity.NewUserID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("FindByID not found: got %v, want ErrNotFound", err)
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)
	user := entity.NewUser(newTestEmail(t, "find@example.com"), newTestUserName(t, "検索対象"))

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := repo.FindByEmail(ctx, newTestEmail(t, "find@example.com"))
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if got.ID() != user.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), user.ID())
	}
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)

	_, err := repo.FindByEmail(ctx, newTestEmail(t, "notexist@example.com"))
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("FindByEmail not found: got %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)
	user := entity.NewUser(newTestEmail(t, "delete@example.com"), newTestUserName(t, "削除対象"))

	if err := repo.Save(ctx, user); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if err := repo.Delete(ctx, user.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, user.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("FindByID after delete: got %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	repo := postgres.NewUserRepository(tx)

	err := repo.Delete(ctx, entity.NewUserID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Delete not found: got %v, want ErrNotFound", err)
	}
}
