package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestUpdate_Success(t *testing.T) {
	existing := newTestTask(t)
	saveCalled := false
	repo := &mockTaskRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.TaskID) (*entity.Task, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Task) error {
			saveCalled = true
			return nil
		},
	}

	dueDate := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	uc := NewUpdate(repo)
	out, err := uc.Execute(context.Background(), UpdateInput{
		UserID:  entity.NewUserID(),
		TaskID:  existing.ID(),
		Title:   "面接準備",
		Type:    "schedule",
		Status:  "done",
		DueDate: &dueDate,
		Notify:  true,
		Memo:    "新しいメモ",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Task.Title().String() != "面接準備" {
		t.Errorf("Title = %q, want %q", out.Task.Title().String(), "面接準備")
	}
	if out.Task.TaskType().String() != "schedule" {
		t.Errorf("TaskType = %q, want %q", out.Task.TaskType().String(), "schedule")
	}
	if out.Task.Status().String() != "done" {
		t.Errorf("Status = %q, want %q", out.Task.Status().String(), "done")
	}
	if out.Task.DueDate() == nil {
		t.Fatal("DueDate should not be nil")
	}
	if !out.Task.DueDate().Equal(dueDate) {
		t.Errorf("DueDate = %v, want %v", out.Task.DueDate(), dueDate)
	}
	if !out.Task.Notify() {
		t.Error("Notify should be true")
	}
	if out.Task.Memo() != "新しいメモ" {
		t.Errorf("Memo = %q, want %q", out.Task.Memo(), "新しいメモ")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestUpdate_ClearDueDate(t *testing.T) {
	existing := newTestTask(t)
	existing.SetDueDate(time.Now())
	repo := &mockTaskRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.TaskID) (*entity.Task, error) {
			return existing, nil
		},
	}

	uc := NewUpdate(repo)
	out, err := uc.Execute(context.Background(), UpdateInput{
		UserID:  entity.NewUserID(),
		TaskID:  existing.ID(),
		Title:   "ES提出",
		Type:    "deadline",
		Status:  "todo",
		DueDate: nil,
		Memo:    "",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Task.DueDate() != nil {
		t.Errorf("DueDate = %v, want nil", out.Task.DueDate())
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID: entity.NewUserID(),
		TaskID: entity.NewTaskID(),
		Title:  "ES提出",
		Type:   "deadline",
		Status: "todo",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestUpdate_InvalidTitle(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID: entity.NewUserID(),
		TaskID: entity.NewTaskID(),
		Title:  "",
		Type:   "deadline",
		Status: "todo",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrTaskTitleEmpty) {
		t.Errorf("error = %v, want ErrTaskTitleEmpty", err)
	}
}

func TestUpdate_InvalidStatus(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID: entity.NewUserID(),
		TaskID: entity.NewTaskID(),
		Title:  "ES提出",
		Type:   "deadline",
		Status: "invalid",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrTaskStatusInvalid) {
		t.Errorf("error = %v, want ErrTaskStatusInvalid", err)
	}
}

func TestUpdate_SaveError(t *testing.T) {
	existing := newTestTask(t)
	saveErr := errors.New("db write failed")
	repo := &mockTaskRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.TaskID) (*entity.Task, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Task) error {
			return saveErr
		},
	}

	uc := NewUpdate(repo)
	_, err := uc.Execute(context.Background(), UpdateInput{
		UserID: entity.NewUserID(),
		TaskID: existing.ID(),
		Title:  "ES提出",
		Type:   "deadline",
		Status: "todo",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
