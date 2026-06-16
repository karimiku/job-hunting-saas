package inboxclip

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

func TestDelete_OwnedClip_Removes(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	userID := entity.NewUserID()
	url, _ := value.NewURL("https://example.com/jobs/1")
	title, _ := value.NewInboxClipTitle("a")
	src, _ := value.NewSource("マイナビ")
	guess, _ := value.NewInboxClipGuess("")
	clip := entity.NewInboxClip(userID, url, title, src, guess, value.InboxClipContentText{})
	_ = repo.Create(context.Background(), clip)

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{UserID: userID, ClipID: clip.ID()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := repo.FindByID(context.Background(), userID, clip.ID()); !errors.Is(err, repository.ErrNotFound) {
		t.Error("clip should be deleted")
	}
}

func TestDelete_OtherUserClip_NotFound(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	owner := entity.NewUserID()
	other := entity.NewUserID()
	url, _ := value.NewURL("https://example.com/jobs/1")
	title, _ := value.NewInboxClipTitle("a")
	src, _ := value.NewSource("マイナビ")
	guess, _ := value.NewInboxClipGuess("")
	clip := entity.NewInboxClip(owner, url, title, src, guess, value.InboxClipContentText{})
	_ = repo.Create(context.Background(), clip)

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{UserID: other, ClipID: clip.ID()})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}
