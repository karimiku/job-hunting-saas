package entity

import (
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestSource(t *testing.T) value.Source {
	t.Helper()
	s, err := value.NewSource("マイナビ")
	if err != nil {
		t.Fatalf("NewSource failed: %v", err)
	}
	return s
}

func newTestStage(t *testing.T, kind, label string) value.Stage {
	t.Helper()
	s, err := value.NewStage(kind, label)
	if err != nil {
		t.Fatalf("NewStage failed: %v", err)
	}
	return s
}

func newTestEntryStatus(t *testing.T, raw string) value.EntryStatus {
	t.Helper()
	s, err := value.NewEntryStatus(raw)
	if err != nil {
		t.Fatalf("NewEntryStatus failed: %v", err)
	}
	return s
}

func TestNewEntry(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	source := newTestSource(t)

	t.Run("valid entry", func(t *testing.T) {
		entry, err := NewEntry(userID, companyID, "本選考", source)
		if err != nil {
			t.Fatalf("NewEntry should succeed, but got error: %v", err)
		}
		if entry.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if entry.UserID() != userID {
			t.Errorf("UserID() = %v, want %v", entry.UserID(), userID)
		}
		if entry.CompanyID() != companyID {
			t.Errorf("CompanyID() = %v, want %v", entry.CompanyID(), companyID)
		}
		if entry.Route() != "本選考" {
			t.Errorf("Route() = %q, want %q", entry.Route(), "本選考")
		}
		if entry.Source().String() != "マイナビ" {
			t.Errorf("Source() = %q, want %q", entry.Source().String(), "マイナビ")
		}
		if entry.Status().String() != "in_progress" {
			t.Errorf("Status() should be in_progress, got %q", entry.Status().String())
		}
		if entry.Stage().Kind() != "application" {
			t.Errorf("Stage().Kind() should be application, got %q", entry.Stage().Kind())
		}
		if entry.Memo() != "" {
			t.Errorf("Memo() should be empty, got %q", entry.Memo())
		}
	})

	t.Run("empty route", func(t *testing.T) {
		_, err := NewEntry(userID, companyID, "", source)
		if err == nil {
			t.Error("NewEntry with empty route should return error")
		}
	})

	t.Run("whitespace route", func(t *testing.T) {
		_, err := NewEntry(userID, companyID, "   ", source)
		if err == nil {
			t.Error("NewEntry with whitespace route should return error")
		}
	})
}

func TestEntry_UpdateStage(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	source := newTestSource(t)
	entry, _ := NewEntry(userID, companyID, "本選考", source)

	newStage := newTestStage(t, "interview", "一次面接")
	entry.UpdateStage(newStage)

	if entry.Stage().Kind() != "interview" {
		t.Errorf("Stage().Kind() = %q, want %q", entry.Stage().Kind(), "interview")
	}
	if entry.Stage().Label() != "一次面接" {
		t.Errorf("Stage().Label() = %q, want %q", entry.Stage().Label(), "一次面接")
	}
}

func TestEntry_UpdateStatus(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	source := newTestSource(t)
	entry, _ := NewEntry(userID, companyID, "本選考", source)

	offered := newTestEntryStatus(t, "offered")
	entry.UpdateStatus(offered)

	if entry.Status().String() != "offered" {
		t.Errorf("Status() = %q, want %q", entry.Status().String(), "offered")
	}
}

func TestEntry_UpdateSource(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	source := newTestSource(t)
	entry, _ := NewEntry(userID, companyID, "本選考", source)

	newSource, err := value.NewSource("リクナビ")
	if err != nil {
		t.Fatalf("NewSource failed: %v", err)
	}
	entry.UpdateSource(newSource)

	if entry.Source().String() != "リクナビ" {
		t.Errorf("Source() = %q, want %q", entry.Source().String(), "リクナビ")
	}
}

func TestEntry_UpdateMemo(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	source := newTestSource(t)
	entry, _ := NewEntry(userID, companyID, "本選考", source)

	entry.UpdateMemo("第一志望")
	if entry.Memo() != "第一志望" {
		t.Errorf("Memo() = %q, want %q", entry.Memo(), "第一志望")
	}
}
