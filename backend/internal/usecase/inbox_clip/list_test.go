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
	url, _ := value.NewURL("https://example.com/jobs/1")
	src, _ := value.NewSource("マイナビ")

	_ = repo.Create(context.Background(), entity.NewInboxClip(owner, url, "a", src, ""))
	_ = repo.Create(context.Background(), entity.NewInboxClip(owner, url, "b", src, ""))
	_ = repo.Create(context.Background(), entity.NewInboxClip(other, url, "c", src, ""))

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{UserID: owner})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Clips) != 2 {
		t.Errorf("len = %d, want 2 (other user excluded)", len(out.Clips))
	}
}
