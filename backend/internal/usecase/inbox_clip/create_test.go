package inboxclip

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
)

func TestCreate_Success(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)
	userID := entity.NewUserID()

	out, err := uc.Execute(context.Background(), CreateInput{
		UserID: userID,
		URL:    "https://job.mynavi.jp/26/pc/search/corp123/outline.html",
		Title:  "○○商事 / 総合職",
		Source: "マイナビ",
		Guess:  "○○商事",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Clip.UserID() != userID {
		t.Errorf("UserID = %v, want %v", out.Clip.UserID(), userID)
	}
	if out.Clip.Source().String() != "マイナビ" {
		t.Errorf("Source = %q", out.Clip.Source().String())
	}
}

func TestCreate_InvalidURL(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		URL:    "javascript:alert(1)",
		Title:  "x",
		Source: "マイナビ",
	})
	if !errors.Is(err, value.ErrURLInvalid) {
		t.Errorf("error = %v, want ErrURLInvalid", err)
	}
}

func TestCreate_EmptySource(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		URL:    "https://example.com/jobs/1",
		Title:  "x",
		Source: "",
	})
	if !errors.Is(err, value.ErrSourceEmpty) {
		t.Errorf("error = %v, want ErrSourceEmpty", err)
	}
}
