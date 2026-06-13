package entity

import (
	"testing"
	"time"
)

func TestNewAIAccessToken(t *testing.T) {
	t.Parallel()

	userID := NewUserID()
	token := NewAIAccessToken(userID, "Codex", "hash", "entre_ai_abcd...wxyz")

	if token.ID().IsZero() {
		t.Error("ID() is zero")
	}
	if token.UserID() != userID {
		t.Errorf("UserID() = %v, want %v", token.UserID(), userID)
	}
	if token.TokenHash() != "hash" {
		t.Errorf("TokenHash() = %q, want hash", token.TokenHash())
	}
	if token.LastUsedAt() != nil {
		t.Error("LastUsedAt() should be nil")
	}
	if token.RevokedAt() != nil {
		t.Error("RevokedAt() should be nil")
	}
}

func TestAIAccessTokenMarkUsedAndRevoke(t *testing.T) {
	t.Parallel()

	token := NewAIAccessToken(NewUserID(), "Codex", "hash", "preview")
	usedAt := time.Date(2026, 6, 13, 10, 0, 0, 0, time.UTC)
	token.MarkUsed(usedAt)

	if token.LastUsedAt() == nil || !token.LastUsedAt().Equal(usedAt) {
		t.Fatalf("LastUsedAt() = %v, want %v", token.LastUsedAt(), usedAt)
	}

	revokedAt := usedAt.Add(time.Hour)
	token.Revoke(revokedAt)
	if token.RevokedAt() == nil || !token.RevokedAt().Equal(revokedAt) {
		t.Fatalf("RevokedAt() = %v, want %v", token.RevokedAt(), revokedAt)
	}
}
