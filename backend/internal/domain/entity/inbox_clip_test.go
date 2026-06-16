package entity

import (
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestNewInboxClip_BasicFields(t *testing.T) {
	userID := NewUserID()
	url, _ := value.NewURL("https://job.mynavi.jp/26/pc/search/corp123/outline.html")
	title, _ := value.NewInboxClipTitle("○○商事 / 総合職")
	source, _ := value.NewSource("マイナビ")
	guess, _ := value.NewInboxClipGuess("○○商事")

	clip := NewInboxClip(userID, url, title, source, guess, value.InboxClipContentText{})

	if clip.UserID() != userID {
		t.Errorf("UserID = %v, want %v", clip.UserID(), userID)
	}
	if clip.URL().String() != "https://job.mynavi.jp/26/pc/search/corp123/outline.html" {
		t.Errorf("URL = %q", clip.URL().String())
	}
	if clip.Title().String() != "○○商事 / 総合職" {
		t.Errorf("Title = %q", clip.Title().String())
	}
	if clip.Source().String() != "マイナビ" {
		t.Errorf("Source = %q", clip.Source().String())
	}
	if clip.Guess().String() != "○○商事" {
		t.Errorf("Guess = %q", clip.Guess().String())
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
	title, _ := value.NewInboxClipTitle("title")
	source, _ := value.NewSource("マイナビ")
	guess, _ := value.NewInboxClipGuess("")

	clip := NewInboxClip(userID, url, title, source, guess, value.InboxClipContentText{})
	if clip.Guess().String() != "" {
		t.Errorf("Guess = %q, want empty", clip.Guess().String())
	}
}

func TestReconstructInboxClip(t *testing.T) {
	id := NewInboxClipID()
	userID := NewUserID()
	url, _ := value.NewURL("https://example.com/jobs/1")
	title, _ := value.NewInboxClipTitle("title")
	source, _ := value.NewSource("リクナビ")
	guess, _ := value.NewInboxClipGuess("guess")

	captured, _ := time.Parse(time.RFC3339, "2026-04-26T00:00:00Z")
	contentText, _ := value.NewInboxClipContentText("選考フロー: ES提出、一次面接")
	clip := ReconstructInboxClip(id, userID, url, title, source, guess, contentText, captured)
	if clip.ID() != id {
		t.Errorf("ID = %v, want %v", clip.ID(), id)
	}
	if clip.UserID() != userID {
		t.Errorf("UserID = %v, want %v", clip.UserID(), userID)
	}
	if clip.ContentText().String() != "選考フロー: ES提出、一次面接" {
		t.Errorf("ContentText = %q", clip.ContentText().String())
	}
}
