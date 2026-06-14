package user

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

type mockAccountRepo struct {
	deleteFn func(ctx context.Context, userID entity.UserID) error
}

func (m *mockAccountRepo) DeleteAccount(ctx context.Context, userID entity.UserID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID)
	}
	return nil
}

func TestDeleteAccount_Success(t *testing.T) {
	userID := entity.NewUserID()
	called := false
	uc := NewDeleteAccount(&mockAccountRepo{
		deleteFn: func(_ context.Context, got entity.UserID) error {
			called = true
			if got != userID {
				t.Errorf("userID = %v, want %v", got, userID)
			}
			return nil
		},
	})

	if err := uc.Execute(context.Background(), DeleteAccountInput{UserID: userID}); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !called {
		t.Fatal("repository should be called")
	}
}

func TestDeleteAccount_EmptyUserID(t *testing.T) {
	called := false
	uc := NewDeleteAccount(&mockAccountRepo{
		deleteFn: func(_ context.Context, _ entity.UserID) error {
			called = true
			return nil
		},
	})

	if err := uc.Execute(context.Background(), DeleteAccountInput{}); err == nil {
		t.Fatal("Execute should return validation error")
	}
	if called {
		t.Fatal("repository should not be called")
	}
}

func TestDeleteAccount_RepositoryError(t *testing.T) {
	expected := errors.New("db failed")
	uc := NewDeleteAccount(&mockAccountRepo{
		deleteFn: func(_ context.Context, _ entity.UserID) error {
			return expected
		},
	})

	err := uc.Execute(context.Background(), DeleteAccountInput{UserID: entity.NewUserID()})
	if !errors.Is(err, expected) {
		t.Fatalf("err = %v, want %v", err, expected)
	}
}
