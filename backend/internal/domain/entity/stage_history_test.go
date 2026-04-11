package entity

import (
	"testing"
)

func TestNewStageHistory(t *testing.T) {
	entryID := NewEntryID()
	stage := newTestStage(t, "interview", "一次面接")

	t.Run("valid stage history", func(t *testing.T) {
		h := NewStageHistory(entryID, stage, "オンライン面接")
		if h.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if h.EntryID() != entryID {
			t.Errorf("EntryID() = %v, want %v", h.EntryID(), entryID)
		}
		if h.Stage().Kind().String() != "interview" {
			t.Errorf("Stage().Kind() = %q, want %q", h.Stage().Kind().String(), "interview")
		}
		if h.Stage().Label() != "一次面接" {
			t.Errorf("Stage().Label() = %q, want %q", h.Stage().Label(), "一次面接")
		}
		if h.Note() != "オンライン面接" {
			t.Errorf("Note() = %q, want %q", h.Note(), "オンライン面接")
		}
	})

	t.Run("empty note", func(t *testing.T) {
		h := NewStageHistory(entryID, stage, "")
		if h.Note() != "" {
			t.Errorf("Note() = %q, want empty", h.Note())
		}
	})
}
