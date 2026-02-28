package entity

import (
	"testing"

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

func newTestUserName(t *testing.T, raw string) value.UserName {
	t.Helper()
	name, err := value.NewUserName(raw)
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	return name
}

func TestNewUser(t *testing.T) {
	email := newTestEmail(t, "test@example.com")
	name := newTestUserName(t, "田中太郎")

	t.Run("valid user", func(t *testing.T) {
		user := NewUser(email, name)
		if user.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if user.Email().String() != "test@example.com" {
			t.Errorf("Email() = %q, want %q", user.Email().String(), "test@example.com")
		}
		if user.Name().String() != "田中太郎" {
			t.Errorf("Name() = %q, want %q", user.Name().String(), "田中太郎")
		}
	})
}

func TestUser_Rename(t *testing.T) {
	email := newTestEmail(t, "test@example.com")
	name := newTestUserName(t, "田中太郎")
	user := NewUser(email, name)

	newName := newTestUserName(t, "佐藤花子")
	user.Rename(newName)

	if user.Name().String() != "佐藤花子" {
		t.Errorf("Name() = %q, want %q", user.Name().String(), "佐藤花子")
	}
}

func TestUser_ChangeEmail(t *testing.T) {
	email := newTestEmail(t, "old@example.com")
	name := newTestUserName(t, "田中太郎")
	user := NewUser(email, name)

	newEmail := newTestEmail(t, "new@example.com")
	user.ChangeEmail(newEmail)

	if user.Email().String() != "new@example.com" {
		t.Errorf("Email() = %q, want %q", user.Email().String(), "new@example.com")
	}
}
