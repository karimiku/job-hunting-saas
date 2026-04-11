package entity

import (
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestNewExternalIdentity(t *testing.T) {
	userID := NewUserID()
	provider := value.AuthProviderGoogle()
	subject := "google-oauth-subject-123"

	identity := NewExternalIdentity(userID, provider, subject)

	if identity.ID().IsZero() {
		t.Error("ID should not be zero")
	}
	if identity.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", identity.UserID(), userID)
	}
	if !identity.Provider().Equals(value.AuthProviderGoogle()) {
		t.Errorf("Provider() = %q, want %q", identity.Provider().String(), "google")
	}
	if identity.Subject() != subject {
		t.Errorf("Subject() = %q, want %q", identity.Subject(), subject)
	}
	if identity.CreatedAt().IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}
