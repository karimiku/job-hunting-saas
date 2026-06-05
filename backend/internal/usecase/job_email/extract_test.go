package jobemail

import (
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.Local)
	got := NewExtract().Execute(ExtractInput{
		Subject: "株式会社テスト 一次面接のご案内",
		Text:    "一次面接は6月12日 14:00です。ES提出期限は6月10日 23:59までです。",
		Now:     now,
	})

	if got.CompanyName != "株式会社テスト" {
		t.Fatalf("CompanyName = %q, want 株式会社テスト", got.CompanyName)
	}
	if got.StageKind != "interview" || got.StageLabel != "一次面接" {
		t.Fatalf("stage = %s/%s, want interview/一次面接", got.StageKind, got.StageLabel)
	}
	if got.EventAt == nil || got.DeadlineAt == nil {
		t.Fatalf("EventAt/DeadlineAt should be extracted: %+v", got)
	}
	if len(got.SuggestedTasks) != 2 {
		t.Fatalf("SuggestedTasks len = %d, want 2", len(got.SuggestedTasks))
	}
}

func TestExtractUsesProvidedCompanyName(t *testing.T) {
	now := time.Date(2026, 6, 4, 12, 0, 0, 0, time.Local)
	got := NewExtract().Execute(ExtractInput{
		Subject:     "面接のご案内",
		Text:        "6月12日 14:00に面接です。",
		CompanyName: "指定会社",
		Now:         now,
	})

	if got.CompanyName != "指定会社" {
		t.Fatalf("CompanyName = %q, want 指定会社", got.CompanyName)
	}
}
