package inboxclip

import (
	"context"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

func TestList_OnlyOwnedReturned(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	owner := entity.NewUserID()
	other := entity.NewUserID()
	urlA, _ := value.NewURL("https://example.com/jobs/1")
	urlB, _ := value.NewURL("https://example.com/jobs/2")
	urlC, _ := value.NewURL("https://example.com/jobs/3")
	src, _ := value.NewSource("マイナビ")
	guess, _ := value.NewInboxClipGuess("")

	titleA, _ := value.NewInboxClipTitle("a")
	titleB, _ := value.NewInboxClipTitle("b")
	titleC, _ := value.NewInboxClipTitle("c")
	_ = repo.Create(context.Background(), entity.NewInboxClip(owner, urlA, titleA, src, guess, value.InboxClipContentText{}))
	_ = repo.Create(context.Background(), entity.NewInboxClip(owner, urlB, titleB, src, guess, value.InboxClipContentText{}))
	_ = repo.Create(context.Background(), entity.NewInboxClip(other, urlC, titleC, src, guess, value.InboxClipContentText{}))

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{UserID: owner})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Clips) != 2 {
		t.Errorf("len = %d, want 2 (other user excluded)", len(out.Clips))
	}
}
