package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestEmail(t *testing.T, raw string) value.Email {
	t.Helper()
	email, err := value.NewEmail(raw)
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	return email
}

func TestNewUser(t *testing.T) {
	email := newTestEmail(t, "test@example.com")

	t.Run("valid user", func(t *testing.T) {
		user, err := NewUser(email, "田中太郎")
		if err != nil {
			t.Fatalf("NewUser should succeed, but got error: %v", err)
		}
		if user.ID() == uuid.Nil {
			t.Error("ID should not be nil")
		}
		if user.Email().String() != "test@example.com" {
			t.Errorf("Email() = %q, want %q", user.Email().String(), "test@example.com")
		}
		if user.Name() != "田中太郎" {
			t.Errorf("Name() = %q, want %q", user.Name(), "田中太郎")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := NewUser(email, "")
		if err == nil {
			t.Error("NewUser with empty name should return error")
		}
	})

	t.Run("whitespace name", func(t *testing.T) {
		_, err := NewUser(email, "   ")
		if err == nil {
			t.Error("NewUser with whitespace name should return error")
		}
	})
}

func TestUser_Rename(t *testing.T) {
	email := newTestEmail(t, "test@example.com")
	user, _ := NewUser(email, "田中太郎")

	t.Run("valid rename", func(t *testing.T) {
		err := user.Rename("佐藤花子")
		if err != nil {
			t.Fatalf("Rename should succeed, but got error: %v", err)
		}
		if user.Name() != "佐藤花子" {
			t.Errorf("Name() = %q, want %q", user.Name(), "佐藤花子")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		err := user.Rename("")
		if err == nil {
			t.Error("Rename with empty name should return error")
		}
	})
}

func TestUser_ChangeEmail(t *testing.T) {
	email := newTestEmail(t, "old@example.com")
	user, _ := NewUser(email, "田中太郎")

	newEmail := newTestEmail(t, "new@example.com")
	user.ChangeEmail(newEmail)

	if user.Email().String() != "new@example.com" {
		t.Errorf("Email() = %q, want %q", user.Email().String(), "new@example.com")
	}
}
