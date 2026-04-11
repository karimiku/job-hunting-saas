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

func entryFound() *mockEntryRepo {
	return &mockEntryRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) (*entity.Entry, error) {
			return &entity.Entry{}, nil
		},
	}
}

func TestCreate_Success(t *testing.T) {
	saveCalled := false
	taskRepo := &mockTaskRepo{
		saveFn: func(_ context.Context, _ *entity.Task) error {
			saveCalled = true
			return nil
		},
	}

	dueDate := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	uc := NewCreate(taskRepo, entryFound())
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "ES提出",
		Type:    "deadline",
		DueDate: &dueDate,
		Memo:    "メモ内容",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Task == nil {
		t.Fatal("Task should not be nil")
	}
	if out.Task.Title().String() != "ES提出" {
		t.Errorf("Title = %q, want %q", out.Task.Title().String(), "ES提出")
	}
	if out.Task.TaskType().String() != "deadline" {
		t.Errorf("TaskType = %q, want %q", out.Task.TaskType().String(), "deadline")
	}
	if out.Task.DueDate() == nil {
		t.Fatal("DueDate should not be nil")
	}
	if !out.Task.DueDate().Equal(dueDate) {
		t.Errorf("DueDate = %v, want %v", out.Task.DueDate(), dueDate)
	}
	if out.Task.Memo() != "メモ内容" {
		t.Errorf("Memo = %q, want %q", out.Task.Memo(), "メモ内容")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestCreate_WithoutOptionalFields(t *testing.T) {
	uc := NewCreate(&mockTaskRepo{}, entryFound())
	out, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "面接準備",
		Type:    "schedule",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Task.DueDate() != nil {
		t.Errorf("DueDate = %v, want nil", out.Task.DueDate())
	}
	if out.Task.Memo() != "" {
		t.Errorf("Memo = %q, want empty", out.Task.Memo())
	}
}

func TestCreate_EntryNotFound(t *testing.T) {
	uc := NewCreate(&mockTaskRepo{}, &mockEntryRepo{})

	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "ES提出",
		Type:    "deadline",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	uc := NewCreate(&mockTaskRepo{}, entryFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "",
		Type:    "deadline",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrTaskTitleEmpty) {
		t.Errorf("error = %v, want ErrTaskTitleEmpty", err)
	}
}

func TestCreate_InvalidType(t *testing.T) {
	uc := NewCreate(&mockTaskRepo{}, entryFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "ES提出",
		Type:    "invalid",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrTaskTypeInvalid) {
		t.Errorf("error = %v, want ErrTaskTypeInvalid", err)
	}
}

func TestCreate_SaveError(t *testing.T) {
	saveErr := errors.New("db write failed")
	taskRepo := &mockTaskRepo{
		saveFn: func(_ context.Context, _ *entity.Task) error {
			return saveErr
		},
	}

	uc := NewCreate(taskRepo, entryFound())
	_, err := uc.Execute(context.Background(), CreateInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
		Title:   "ES提出",
		Type:    "deadline",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
