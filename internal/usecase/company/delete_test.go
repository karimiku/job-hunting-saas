package company

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestDelete_Success(t *testing.T) {
	deleteCalled := false
	repo := &mockCompanyRepo{
		deleteFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) error {
			deleteCalled = true
			return nil
		},
	}

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
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
	repo := &mockCompanyRepo{
		deleteFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) error {
			return dbErr
		},
	}

	uc := NewDelete(repo)
	err := uc.Execute(context.Background(), DeleteInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
