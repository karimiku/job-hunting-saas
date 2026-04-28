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

func mustParseTime(t *testing.T, raw string) time.Time {
	t.Helper()
	tt, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		t.Fatalf("time.Parse(%q): %v", raw, err)
	}
	return tt
}

func newTestURL(t *testing.T, raw string) value.URL {
	t.Helper()
	u, err := value.NewURL(raw)
	if err != nil {
		t.Fatalf("NewURL(%q) failed: %v", raw, err)
	}
	return u
}

func newTestClipSource(t *testing.T, raw string) value.Source {
	t.Helper()
	s, err := value.NewSource(raw)
	if err != nil {
		t.Fatalf("NewSource(%q) failed: %v", raw, err)
	}
	return s
}

func TestInboxClipRepository_Create_Find(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	clip := entity.NewInboxClip(
		userID,
		newTestURL(t, "https://job.mynavi.jp/26/pc/search/corp123/outline.html"),
		"株式会社サンプル — 募集要項",
		newTestClipSource(t, "マイナビ"),
		"株式会社サンプル",
	)

	if err := repo.Create(ctx, clip); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, clip.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if got.ID() != clip.ID() {
		t.Errorf("ID = %v, want %v", got.ID(), clip.ID())
	}
	if got.UserID() != userID {
		t.Errorf("UserID = %v, want %v", got.UserID(), userID)
	}
	if got.URL().String() != clip.URL().String() {
		t.Errorf("URL = %q, want %q", got.URL().String(), clip.URL().String())
	}
	if got.Title() != clip.Title() {
		t.Errorf("Title = %q, want %q", got.Title(), clip.Title())
	}
	if got.Source().String() != "マイナビ" {
		t.Errorf("Source = %q, want %q", got.Source().String(), "マイナビ")
	}
	if got.Guess() != "株式会社サンプル" {
		t.Errorf("Guess = %q, want %q", got.Guess(), "株式会社サンプル")
	}
	if !got.CapturedAt().Equal(clip.CapturedAt()) {
		t.Errorf("CapturedAt = %v, want %v", got.CapturedAt(), clip.CapturedAt())
	}
}

func TestInboxClipRepository_Create_EmptyGuess(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	clip := entity.NewInboxClip(
		userID,
		newTestURL(t, "https://example.com/jobs/1"),
		"求人タイトル",
		newTestClipSource(t, "リクナビ"),
		"",
	)

	if err := repo.Create(ctx, clip); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	got, err := repo.FindByID(ctx, userID, clip.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if got.Guess() != "" {
		t.Errorf("Guess = %q, want empty", got.Guess())
	}
}

func TestInboxClipRepository_FindByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	_, err := repo.FindByID(ctx, userID, entity.NewInboxClipID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestInboxClipRepository_FindByID_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	clip := entity.NewInboxClip(userID, newTestURL(t, "https://example.com/jobs/1"), "T", newTestClipSource(t, "マイナビ"), "")
	if err := repo.Create(ctx, clip); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	_, err := repo.FindByID(ctx, otherUserID, clip.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestInboxClipRepository_ListByUserID_OrderByCapturedAtDesc(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	// time.Now() ベースの NewInboxClip だと連続生成で同時刻になり得るので、Reconstruct で明示的に時刻を分ける。
	older := entity.ReconstructInboxClip(
		entity.NewInboxClipID(),
		userID,
		newTestURL(t, "https://example.com/older"),
		"older",
		newTestClipSource(t, "マイナビ"),
		"",
		mustParseTime(t, "2026-04-01T10:00:00Z"),
	)
	newer := entity.ReconstructInboxClip(
		entity.NewInboxClipID(),
		userID,
		newTestURL(t, "https://example.com/newer"),
		"newer",
		newTestClipSource(t, "マイナビ"),
		"",
		mustParseTime(t, "2026-04-02T10:00:00Z"),
	)
	for _, c := range []*entity.InboxClip{older, newer} {
		if err := repo.Create(ctx, c); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	list, err := repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("len = %d, want 2", len(list))
	}
	if list[0].ID() != newer.ID() {
		t.Errorf("list[0] = %v, want %v (newer first)", list[0].ID(), newer.ID())
	}
	if list[1].ID() != older.ID() {
		t.Errorf("list[1] = %v, want %v (older second)", list[1].ID(), older.ID())
	}
}

func TestInboxClipRepository_ListByUserID_FiltersByOwner(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	mine := entity.NewInboxClip(userID, newTestURL(t, "https://example.com/mine"), "mine", newTestClipSource(t, "マイナビ"), "")
	theirs := entity.NewInboxClip(otherUserID, newTestURL(t, "https://example.com/theirs"), "theirs", newTestClipSource(t, "マイナビ"), "")
	for _, c := range []*entity.InboxClip{mine, theirs} {
		if err := repo.Create(ctx, c); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	list, err := repo.ListByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUserID failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ID() != mine.ID() {
		t.Errorf("ID = %v, want %v", list[0].ID(), mine.ID())
	}
}

func TestInboxClipRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	clip := entity.NewInboxClip(userID, newTestURL(t, "https://example.com/del"), "T", newTestClipSource(t, "マイナビ"), "")
	if err := repo.Create(ctx, clip); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(ctx, userID, clip.ID()); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(ctx, userID, clip.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("after delete: err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestInboxClipRepository_Delete_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	err := repo.Delete(ctx, userID, entity.NewInboxClipID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}
}

func TestInboxClipRepository_Delete_WrongUser(t *testing.T) {
	pool := setupTestDB(t)
	tx := beginTx(t, pool)
	ctx := context.Background()

	userID := insertTestUser(t, tx)
	otherUserID := insertTestUser(t, tx)
	repo := postgres.NewInboxClipRepository(tx)

	clip := entity.NewInboxClip(userID, newTestURL(t, "https://example.com/x"), "T", newTestClipSource(t, "マイナビ"), "")
	if err := repo.Create(ctx, clip); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	err := repo.Delete(ctx, otherUserID, clip.ID())
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("err = %v, want %v", err, repository.ErrNotFound)
	}

	got, err := repo.FindByID(ctx, userID, clip.ID())
	if err != nil {
		t.Fatalf("FindByID after wrong-user delete: %v", err)
	}
	if got.ID() != clip.ID() {
		t.Errorf("clip should still exist, got ID %v", got.ID())
	}
}
