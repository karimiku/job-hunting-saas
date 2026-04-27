package entity

import (
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestNewInboxClip_BasicFields(t *testing.T) {
	userID := NewUserID()
	url, _ := value.NewURL("https://job.mynavi.jp/26/pc/search/corp123/outline.html")
	source, _ := value.NewSource("マイナビ")

	clip := NewInboxClip(userID, url, "○○商事 / 総合職", source, "○○商事")

	if clip.UserID() != userID {
		t.Errorf("UserID = %v, want %v", clip.UserID(), userID)
	}
	if clip.URL().String() != "https://job.mynavi.jp/26/pc/search/corp123/outline.html" {
		t.Errorf("URL = %q", clip.URL().String())
	}
	if clip.Title() != "○○商事 / 総合職" {
		t.Errorf("Title = %q", clip.Title())
	}
	if clip.Source().String() != "マイナビ" {
		t.Errorf("Source = %q", clip.Source().String())
	}
	if clip.Guess() != "○○商事" {
		t.Errorf("Guess = %q", clip.Guess())
	}
	if clip.ID().IsZero() {
		t.Error("ID should be generated")
	}
	if clip.CapturedAt().IsZero() {
		t.Error("CapturedAt should be set")
	}
}

func TestNewInboxClip_GuessIsOptional(t *testing.T) {
	userID := NewUserID()
	url, _ := value.NewURL("https://job.mynavi.jp/26/pc/search/corp123/outline.html")
	source, _ := value.NewSource("マイナビ")

	clip := NewInboxClip(userID, url, "title", source, "")
	if clip.Guess() != "" {
		t.Errorf("Guess = %q, want empty", clip.Guess())
	}
}

func TestReconstructInboxClip(t *testing.T) {
	id := NewInboxClipID()
	userID := NewUserID()
	url, _ := value.NewURL("https://example.com/jobs/1")
	source, _ := value.NewSource("リクナビ")

	captured, _ := time.Parse(time.RFC3339, "2026-04-26T00:00:00Z")
	clip := ReconstructInboxClip(id, userID, url, "title", source, "guess", captured)
	if clip.ID() != id {
		t.Errorf("ID = %v, want %v", clip.ID(), id)
	}
	if clip.UserID() != userID {
		t.Errorf("UserID = %v, want %v", clip.UserID(), userID)
	}
}
