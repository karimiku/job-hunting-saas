package inboxclip

import (
	"context"
	"errors"
	"strings"
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

func TestCreate_DuplicateURL_ReturnsExistingWithoutCreatingSecond(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)
	userID := entity.NewUserID()

	const url = "https://job.mynavi.jp/26/pc/search/corp123/outline.html"

	first, err := uc.Execute(context.Background(), CreateInput{
		UserID: userID,
		URL:    url,
		Title:  "○○商事 / 総合職",
		Source: "マイナビ",
		Guess:  "○○商事",
	})
	if err != nil {
		t.Fatalf("first create: unexpected error: %v", err)
	}

	// 同じ URL をもう一度保存しても、新規作成せず既存クリップを返す（冪等）。
	second, err := uc.Execute(context.Background(), CreateInput{
		UserID: userID,
		URL:    url,
		Title:  "別タイトルで再保存",
		Source: "リクナビ",
		Guess:  "△△",
	})
	if err != nil {
		t.Fatalf("second create: unexpected error: %v", err)
	}

	if second.Clip.ID() != first.Clip.ID() {
		t.Errorf("second clip ID = %v, want existing %v", second.Clip.ID(), first.Clip.ID())
	}

	clips, err := repo.ListByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("list: unexpected error: %v", err)
	}
	if len(clips) != 1 {
		t.Errorf("clip count = %d, want 1 (duplicate URL must not create a second clip)", len(clips))
	}
}

func TestCreate_SameURLDifferentUser_CreatesSeparateClip(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)
	userA := entity.NewUserID()
	userB := entity.NewUserID()

	const url = "https://example.com/jobs/1"

	if _, err := uc.Execute(context.Background(), CreateInput{
		UserID: userA, URL: url, Title: "a", Source: "マイナビ",
	}); err != nil {
		t.Fatalf("userA create: %v", err)
	}
	if _, err := uc.Execute(context.Background(), CreateInput{
		UserID: userB, URL: url, Title: "b", Source: "マイナビ",
	}); err != nil {
		t.Fatalf("userB create: %v", err)
	}

	// 重複抑止はユーザー単位。別ユーザーの同一 URL は独立して作成される。
	for _, u := range []entity.UserID{userA, userB} {
		clips, err := repo.ListByUserID(context.Background(), u)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(clips) != 1 {
			t.Errorf("user %v clip count = %d, want 1", u, len(clips))
		}
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

func TestCreate_EmptyTitle(t *testing.T) {
	repo := inmemory.NewInboxClipRepository()
	uc := NewCreate(repo)

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID: entity.NewUserID(),
		URL:    "https://example.com/jobs/1",
		Title:  "",
		Source: "マイナビ",
	})
	if !errors.Is(err, value.ErrInboxClipTitleEmpty) {
		t.Errorf("error = %v, want ErrInboxClipTitleEmpty", err)
	}
}

func TestCreate_TitleLengthBoundary(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr error
	}{
		{"max length ok", strings.Repeat("あ", value.InboxClipTitleMaxLength), nil},
		{"too long", strings.Repeat("あ", value.InboxClipTitleMaxLength+1), value.ErrInboxClipTitleTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewInboxClipRepository()
			uc := NewCreate(repo)
			_, err := uc.Execute(context.Background(), CreateInput{
				UserID: entity.NewUserID(),
				URL:    "https://example.com/jobs/1",
				Title:  tt.title,
				Source: "マイナビ",
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreate_GuessLengthBoundary(t *testing.T) {
	tests := []struct {
		name    string
		guess   string
		wantErr error
	}{
		{"max length ok", strings.Repeat("あ", value.InboxClipGuessMaxLength), nil},
		{"too long", strings.Repeat("あ", value.InboxClipGuessMaxLength+1), value.ErrInboxClipGuessTooLong},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := inmemory.NewInboxClipRepository()
			uc := NewCreate(repo)
			_, err := uc.Execute(context.Background(), CreateInput{
				UserID: entity.NewUserID(),
				URL:    "https://example.com/jobs/1",
				Title:  "x",
				Source: "マイナビ",
				Guess:  tt.guess,
			})
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
