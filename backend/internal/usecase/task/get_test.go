package task

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestTask(t *testing.T) *entity.Task {
	t.Helper()
	title, err := value.NewTaskTitle("ES提出")
	if err != nil {
		t.Fatalf("NewTaskTitle failed: %v", err)
	}
	taskType, err := value.NewTaskType("deadline")
	if err != nil {
		t.Fatalf("NewTaskType failed: %v", err)
	}
	return entity.NewTask(entity.NewEntryID(), title, taskType)
}

func TestGet_Found(t *testing.T) {
	expected := newTestTask(t)
	repo := &mockTaskRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.TaskID) (*entity.Task, error) {
			return expected, nil
		},
	}

	uc := NewGet(repo)
	out, err := uc.Execute(context.Background(), GetInput{
		UserID: entity.NewUserID(),
		TaskID: expected.ID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Task != expected {
		t.Error("Task should be the expected task")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID: entity.NewUserID(),
		TaskID: entity.NewTaskID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestGet_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockTaskRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.TaskID) (*entity.Task, error) {
			return nil, dbErr
		},
	}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID: entity.NewUserID(),
		TaskID: entity.NewTaskID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
