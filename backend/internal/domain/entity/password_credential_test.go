package entity

import (
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestPassword(t *testing.T, raw string) value.Password {
	t.Helper()
	pw, err := value.NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword failed: %v", err)
	}
	return pw
}

func TestNewPasswordCredential(t *testing.T) {
	userID := NewUserID()
	pw := newTestPassword(t, "my-secret-password")

	cred := NewPasswordCredential(userID, pw)

	if cred.ID().IsZero() {
		t.Error("ID should not be zero")
	}
	if cred.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", cred.UserID(), userID)
	}
	if !cred.Password().Verify("my-secret-password") {
		t.Error("Password should verify correctly")
	}
	if cred.CreatedAt().IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if cred.UpdatedAt().IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestPasswordCredential_ChangePassword(t *testing.T) {
	userID := NewUserID()
	pw := newTestPassword(t, "old-password")
	cred := NewPasswordCredential(userID, pw)

	oldUpdatedAt := cred.UpdatedAt()

	newPw := newTestPassword(t, "new-password")
	cred.ChangePassword(newPw)

	if cred.Password().Verify("old-password") {
		t.Error("Old password should no longer verify")
	}
	if !cred.Password().Verify("new-password") {
		t.Error("New password should verify correctly")
	}
	if !cred.UpdatedAt().After(oldUpdatedAt) || cred.UpdatedAt().Equal(oldUpdatedAt) {
		t.Error("UpdatedAt should be updated after ChangePassword")
	}
}
