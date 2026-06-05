package entity

import (
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestESMemoCategory(t *testing.T, raw string) value.ESMemoCategory {
	t.Helper()
	v, err := value.NewESMemoCategory(raw)
	if err != nil {
		t.Fatalf("NewESMemoCategory() failed: %v", err)
	}
	return v
}

func newTestESMemoTitle(t *testing.T, raw string) value.ESMemoTitle {
	t.Helper()
	v, err := value.NewESMemoTitle(raw)
	if err != nil {
		t.Fatalf("NewESMemoTitle() failed: %v", err)
	}
	return v
}

func newTestESMemoContent(t *testing.T, raw string) value.ESMemoContent {
	t.Helper()
	v, err := value.NewESMemoContent(raw)
	if err != nil {
		t.Fatalf("NewESMemoContent() failed: %v", err)
	}
	return v
}

func newTestESMemoSource(t *testing.T, raw string) value.ESMemoSource {
	t.Helper()
	v, err := value.NewESMemoSource(raw)
	if err != nil {
		t.Fatalf("NewESMemoSource() failed: %v", err)
	}
	return v
}

func TestNewESMemo(t *testing.T) {
	userID := NewUserID()
	entryID := NewEntryID()
	memo := NewESMemo(
		userID,
		&entryID,
		newTestESMemoCategory(t, "interview"),
		newTestESMemoTitle(t, "面接で話す改善経験"),
		newTestESMemoContent(t, "顧客課題を分解して改善した"),
		newTestESMemoSource(t, "mcp"),
	)

	if memo.ID().IsZero() {
		t.Error("ID should not be zero")
	}
	if memo.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", memo.UserID(), userID)
	}
	if memo.EntryID() == nil || *memo.EntryID() != entryID {
		t.Errorf("EntryID() = %v, want %v", memo.EntryID(), entryID)
	}
	if memo.Category().String() != "interview" {
		t.Errorf("Category() = %q", memo.Category().String())
	}
	if memo.Title().String() != "面接で話す改善経験" {
		t.Errorf("Title() = %q", memo.Title().String())
	}
	if memo.Content().String() != "顧客課題を分解して改善した" {
		t.Errorf("Content() = %q", memo.Content().String())
	}
	if memo.Source().String() != "mcp" {
		t.Errorf("Source() = %q", memo.Source().String())
	}
	if memo.CreatedAt().IsZero() || memo.UpdatedAt().IsZero() {
		t.Error("timestamps should not be zero")
	}
}

func TestNewESMemo_AllowsNilEntry(t *testing.T) {
	memo := NewESMemo(
		NewUserID(),
		nil,
		newTestESMemoCategory(t, "general"),
		newTestESMemoTitle(t, "全体メモ"),
		newTestESMemoContent(t, "応募先に紐づかない経験メモ"),
		newTestESMemoSource(t, "mcp"),
	)

	if memo.EntryID() != nil {
		t.Errorf("EntryID() = %v, want nil", memo.EntryID())
	}
}

func TestReconstructESMemo(t *testing.T) {
	id := NewESMemoID()
	userID := NewUserID()
	createdAt := time.Date(2026, 6, 4, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)

	memo := ReconstructESMemo(
		id,
		userID,
		nil,
		newTestESMemoCategory(t, "gakuchika"),
		newTestESMemoTitle(t, "学生時代に力を入れたこと"),
		newTestESMemoContent(t, "継続的に改善した"),
		newTestESMemoSource(t, "mail"),
		createdAt,
		updatedAt,
	)

	if memo.ID() != id {
		t.Errorf("ID() = %v, want %v", memo.ID(), id)
	}
	if memo.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", memo.UserID(), userID)
	}
	if !memo.CreatedAt().Equal(createdAt) || !memo.UpdatedAt().Equal(updatedAt) {
		t.Errorf("timestamps = %v/%v", memo.CreatedAt(), memo.UpdatedAt())
	}
}
