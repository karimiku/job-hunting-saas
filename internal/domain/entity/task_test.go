package entity

import (
	"testing"
	"time"

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

func newTestTaskTitle(t *testing.T, raw string) value.TaskTitle {
	t.Helper()
	title, err := value.NewTaskTitle(raw)
	if err != nil {
		t.Fatalf("NewTaskTitle failed: %v", err)
	}
	return title
}

func TestNewTask(t *testing.T) {
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")

	t.Run("valid task", func(t *testing.T) {
		task := NewTask(entryID, title, taskType)
		if task.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if task.EntryID() != entryID {
			t.Errorf("EntryID() = %v, want %v", task.EntryID(), entryID)
		}
		if task.Title().String() != "ES提出" {
			t.Errorf("Title() = %q, want %q", task.Title().String(), "ES提出")
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
}

func TestTask_Complete(t *testing.T) {
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")
	task := NewTask(entryID, title, taskType)

	task.Complete()
	if !task.Status().IsDone() {
		t.Error("Status should be done after Complete()")
	}
}

func TestTask_Uncomplete(t *testing.T) {
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")
	task := NewTask(entryID, title, taskType)

	task.Complete()
	task.Uncomplete()
	if task.Status().IsDone() {
		t.Error("Status should be todo after Uncomplete()")
	}
}

func TestTask_SetDueDate(t *testing.T) {
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "schedule")
	title := newTestTaskTitle(t, "一次面接")
	task := NewTask(entryID, title, taskType)

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
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")
	task := NewTask(entryID, title, taskType)

	due := time.Date(2026, 3, 15, 14, 0, 0, 0, time.Local)
	task.SetDueDate(due)
	task.ClearDueDate()

	if task.DueDate() != nil {
		t.Error("DueDate() should be nil after ClearDueDate")
	}
}

func TestTask_SetNotify(t *testing.T) {
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")
	task := NewTask(entryID, title, taskType)

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
	entryID := NewEntryID()
	taskType := newTestTaskType(t, "deadline")
	title := newTestTaskTitle(t, "ES提出")
	task := NewTask(entryID, title, taskType)

	task.UpdateMemo("早めに出す")
	if task.Memo() != "早めに出す" {
		t.Errorf("Memo() = %q, want %q", task.Memo(), "早めに出す")
	}
}
