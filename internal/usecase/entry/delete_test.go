package entry

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestDelete_Success(t *testing.T) {
	deleteCalled := false
	repo := &mockEntryRepo{
		deleteFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) error {
			deleteCalled = true
			return nil
		},
	}

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{
		UserID:  entity.NewUserID(),
		EntryID: entity.NewEntryID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleteCalled {
		t.Error("Delete should be called")
	}
}

func TestDelete_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockEntryRepo{
		deleteFn: func(_ context.Context, _ entity.UserID, _ entity.EntryID) error {
			return dbErr
		},
	}

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{
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
