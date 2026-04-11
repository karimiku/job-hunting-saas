package task

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestList_Multiple(t *testing.T) {
	t1 := newTestTask(t)
	t2 := newTestTask(t)
	repo := &mockTaskRepo{
		listByEntryIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) ([]*entity.Task, error) {
			return []*entity.Task{t1, t2}, nil
		},
	}

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Tasks) != 2 {
		t.Fatalf("len(Tasks) = %d, want 2", len(out.Tasks))
	}
}

func TestList_Empty(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewList(repo)
	out, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Tasks) != 0 {
		t.Fatalf("len(Tasks) = %d, want 0", len(out.Tasks))
	}
}

func TestList_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockTaskRepo{
		listByEntryIDFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) ([]*entity.Task, error) {
			return nil, dbErr
		},
	}

	uc := NewList(repo)
	_, err := uc.Execute(context.Background(), ListInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
