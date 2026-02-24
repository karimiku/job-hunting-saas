package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestTaskType(t *testing.T, raw string) value.TaskType {
	t.Helper()
	tt, err := value.NewTaskType(raw)
	if err != nil {
		t.Fatalf("NewTaskType failed: %v", err)
	}
	return tt
}

func TestNewTask(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")

	t.Run("valid task", func(t *testing.T) {
		task, err := NewTask(entryID, "ES提出", taskType)
		if err != nil {
			t.Fatalf("NewTask should succeed, but got error: %v", err)
		}
		if task.ID() == uuid.Nil {
			t.Error("ID should not be nil")
		}
		if task.EntryID() != entryID {
			t.Errorf("EntryID() = %v, want %v", task.EntryID(), entryID)
		}
		if task.Title() != "ES提出" {
			t.Errorf("Title() = %q, want %q", task.Title(), "ES提出")
		}
		if task.TaskType().String() != "deadline" {
			t.Errorf("TaskType() = %q, want %q", task.TaskType().String(), "deadline")
		}
		if task.Status().String() != "todo" {
			t.Errorf("Status() should be todo, got %q", task.Status().String())
		}
		if task.DueDate() != nil {
			t.Error("DueDate() should be nil initially")
		}
		if task.Notify() != false {
			t.Error("Notify() should be false initially")
		}
		if task.Memo() != "" {
			t.Errorf("Memo() should be empty, got %q", task.Memo())
		}
	})

	t.Run("empty title", func(t *testing.T) {
		_, err := NewTask(entryID, "", taskType)
		if err == nil {
			t.Error("NewTask with empty title should return error")
		}
	})

	t.Run("whitespace title", func(t *testing.T) {
		_, err := NewTask(entryID, "   ", taskType)
		if err == nil {
			t.Error("NewTask with whitespace title should return error")
		}
	})
}

func TestTask_Complete(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")
	task, _ := NewTask(entryID, "ES提出", taskType)

	task.Complete()
	if !task.Status().IsDone() {
		t.Error("Status should be done after Complete()")
	}
}

func TestTask_Uncomplete(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")
	task, _ := NewTask(entryID, "ES提出", taskType)

	task.Complete()
	task.Uncomplete()
	if task.Status().IsDone() {
		t.Error("Status should be todo after Uncomplete()")
	}
}

func TestTask_SetDueDate(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "schedule")
	task, _ := NewTask(entryID, "一次面接", taskType)

	due := time.Date(2026, 3, 15, 14, 0, 0, 0, time.Local)
	task.SetDueDate(due)

	if task.DueDate() == nil {
		t.Fatal("DueDate() should not be nil after SetDueDate")
	}
	if !task.DueDate().Equal(due) {
		t.Errorf("DueDate() = %v, want %v", task.DueDate(), due)
	}
}

func TestTask_ClearDueDate(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")
	task, _ := NewTask(entryID, "ES提出", taskType)

	due := time.Date(2026, 3, 15, 14, 0, 0, 0, time.Local)
	task.SetDueDate(due)
	task.ClearDueDate()

	if task.DueDate() != nil {
		t.Error("DueDate() should be nil after ClearDueDate")
	}
}

func TestTask_SetNotify(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")
	task, _ := NewTask(entryID, "ES提出", taskType)

	task.SetNotify(true)
	if task.Notify() != true {
		t.Error("Notify() should be true")
	}

	task.SetNotify(false)
	if task.Notify() != false {
		t.Error("Notify() should be false")
	}
}

func TestTask_UpdateMemo(t *testing.T) {
	entryID := NewID()
	taskType := newTestTaskType(t, "deadline")
	task, _ := NewTask(entryID, "ES提出", taskType)

	task.UpdateMemo("早めに出す")
	if task.Memo() != "早めに出す" {
		t.Errorf("Memo() = %q, want %q", task.Memo(), "早めに出す")
	}
}
