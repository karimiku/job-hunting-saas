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

func newTestInboxClip(t *testing.T, userID entity.UserID, rawURL string) *entity.InboxClip {
	t.Helper()
	url, err := value.NewURL(rawURL)
	if err != nil {
		t.Fatalf("NewURL(%q): %v", rawURL, err)
	}
	title, err := value.NewInboxClipTitle("求人タイトル")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}
	source, err := value.NewSource("マイナビ")
	if err != nil {
		t.Fatalf("NewSource: %v", err)
	}
	guess, err := value.NewInboxClipGuess("")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}
	return entity.NewInboxClip(userID, url, title, source, guess, value.InboxClipContentText{})
}

func TestInboxClipRepository_Create_DuplicateURLForSameUser(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewInboxClipRepository()
	userID := entity.NewUserID()

	first := newTestInboxClip(t, userID, "https://example.com/jobs/1")
	second := newTestInboxClip(t, userID, "https://example.com/jobs/1")

	if err := repo.Create(ctx, first); err != nil {
		t.Fatalf("first Create failed: %v", err)
	}
	if err := repo.Create(ctx, second); !errors.Is(err, repository.ErrAlreadyExists) {
		t.Fatalf("second Create err = %v, want ErrAlreadyExists", err)
	}
}

func TestInboxClipRepository_Create_SameURLDifferentUser(t *testing.T) {
	ctx := context.Background()
	repo := inmemory.NewInboxClipRepository()

	if err := repo.Create(ctx, newTestInboxClip(t, entity.NewUserID(), "https://example.com/jobs/1")); err != nil {
		t.Fatalf("first Create failed: %v", err)
	}
	if err := repo.Create(ctx, newTestInboxClip(t, entity.NewUserID(), "https://example.com/jobs/1")); err != nil {
		t.Fatalf("second Create failed: %v", err)
	}
}
