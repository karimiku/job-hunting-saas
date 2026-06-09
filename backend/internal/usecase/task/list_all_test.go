package task

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestListAll_Multiple(t *testing.T) {
	t1 := newTestTask(t)
	t2 := newTestTask(t)
	repo := &mockTaskRepo{
		listByUserIDFn: func(_ context.Context, _ entity.UserID) ([]*entity.Task, error) {
			return []*entity.Task{t1, t2}, nil
		},
	}

	uc := NewListAll(repo)
	out, err := uc.Execute(context.Background(), ListAllInput{
		UserID: entity.NewUserID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Tasks) != 2 {
		t.Fatalf("len(Tasks) = %d, want 2", len(out.Tasks))
	}
}

func TestListAll_Empty(t *testing.T) {
	repo := &mockTaskRepo{}

	uc := NewListAll(repo)
	out, err := uc.Execute(context.Background(), ListAllInput{
		UserID: entity.NewUserID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Tasks) != 0 {
		t.Fatalf("len(Tasks) = %d, want 0", len(out.Tasks))
	}
}

func TestListAll_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockTaskRepo{
		listByUserIDFn: func(_ context.Context, _ entity.UserID) ([]*entity.Task, error) {
			return nil, dbErr
		},
	}

	uc := NewListAll(repo)
	_, err := uc.Execute(context.Background(), ListAllInput{
		UserID: entity.NewUserID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
