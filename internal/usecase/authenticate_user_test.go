package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// --- mock ---

type mockUserRepo struct {
	findByEmailFn func(ctx context.Context, email value.Email) (*entity.User, error)
	saveFn        func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepo) Save(ctx context.Context, user *entity.User) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, _ entity.UserID) (*entity.User, error) {
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, repository.ErrNotFound
}

func (m *mockUserRepo) Delete(_ context.Context, _ entity.UserID) error {
	return nil
}

// --- helpers ---

func newTestEmail(t *testing.T, raw string) value.Email {
	t.Helper()
	email, err := value.NewEmail(raw)
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	return email
}

func newTestUserName(t *testing.T, raw string) value.UserName {
	t.Helper()
	name, err := value.NewUserName(raw)
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	return name
}

func newExistingUser(t *testing.T) *entity.User {
	t.Helper()
	return entity.NewUser(
		newTestEmail(t, "existing@example.com"),
		newTestUserName(t, "既存ユーザー"),
	)
}

// --- tests ---

func TestAuthenticateUser_NewUser(t *testing.T) {
	saveCalled := false
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
			return nil, repository.ErrNotFound
		},
		saveFn: func(_ context.Context, _ *entity.User) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewAuthenticateUser(repo)
	out, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "new@example.com",
		Name:  "新規ユーザー",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !out.Created {
		t.Error("Created should be true for new user")
	}
	if out.User == nil {
		t.Fatal("User should not be nil")
	}
	if out.User.Email().String() != "new@example.com" {
		t.Errorf("Email = %q, want %q", out.User.Email().String(), "new@example.com")
	}
	if !saveCalled {
		t.Error("Save should be called for new user")
	}
}

func TestAuthenticateUser_ExistingUser(t *testing.T) {
	existing := newExistingUser(t)
	saveCalled := false
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.User) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewAuthenticateUser(repo)
	out, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "existing@example.com",
		Name:  "既存ユーザー",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Created {
		t.Error("Created should be false for existing user")
	}
	if out.User != existing {
		t.Error("User should be the existing user")
	}
	if saveCalled {
		t.Error("Save should not be called for existing user")
	}
}

func TestAuthenticateUser_EmailValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewAuthenticateUser(repo)

	_, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "",
		Name:  "テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrEmailEmpty) {
		t.Errorf("error = %v, want ErrEmailEmpty", err)
	}
}

func TestAuthenticateUser_UserNameValidationError(t *testing.T) {
	repo := &mockUserRepo{}
	uc := NewAuthenticateUser(repo)

	_, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "test@example.com",
		Name:  "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrUserNameEmpty) {
		t.Errorf("error = %v, want ErrUserNameEmpty", err)
	}
}

func TestAuthenticateUser_FindByEmailError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
			return nil, dbErr
		},
	}

	uc := NewAuthenticateUser(repo)
	_, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "test@example.com",
		Name:  "テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}

func TestAuthenticateUser_SaveError(t *testing.T) {
	saveErr := errors.New("db write failed")
	repo := &mockUserRepo{
		findByEmailFn: func(_ context.Context, _ value.Email) (*entity.User, error) {
			return nil, repository.ErrNotFound
		},
		saveFn: func(_ context.Context, _ *entity.User) error {
			return saveErr
		},
	}

	uc := NewAuthenticateUser(repo)
	_, err := uc.Execute(context.Background(), AuthenticateUserInput{
		Email: "test@example.com",
		Name:  "テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
