package entity

import (
	"testing"
)

func TestNewCompany(t *testing.T) {
	userID := NewUserID()

	t.Run("valid company", func(t *testing.T) {
		company, err := NewCompany(userID, "トヨタ自動車")
		if err != nil {
			t.Fatalf("NewCompany should succeed, but got error: %v", err)
		}
		if company.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if company.UserID() != userID {
			t.Errorf("UserID() = %v, want %v", company.UserID(), userID)
		}
		if company.Name() != "トヨタ自動車" {
			t.Errorf("Name() = %q, want %q", company.Name(), "トヨタ自動車")
		}
		if company.Memo() != "" {
			t.Errorf("Memo() should be empty, got %q", company.Memo())
		}
		if company.CreatedAt().IsZero() {
			t.Error("CreatedAt should not be zero")
		}
		if company.UpdatedAt().IsZero() {
			t.Error("UpdatedAt should not be zero")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := NewCompany(userID, "")
		if err == nil {
			t.Error("NewCompany with empty name should return error")
		}
	})

	t.Run("whitespace only name", func(t *testing.T) {
		_, err := NewCompany(userID, "   ")
		if err == nil {
			t.Error("NewCompany with whitespace name should return error")
		}
	})
}

func TestCompany_Rename(t *testing.T) {
	userID := NewUserID()
	company, _ := NewCompany(userID, "旧社名")

	t.Run("valid rename", func(t *testing.T) {
		err := company.Rename("新社名")
		if err != nil {
			t.Fatalf("Rename should succeed, but got error: %v", err)
		}
		if company.Name() != "新社名" {
			t.Errorf("Name() = %q, want %q", company.Name(), "新社名")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		err := company.Rename("")
		if err == nil {
			t.Error("Rename with empty name should return error")
		}
	})
}

func TestCompany_UpdateMemo(t *testing.T) {
	userID := NewUserID()
	company, _ := NewCompany(userID, "テスト株式会社")

	company.UpdateMemo("いい会社")
	if company.Memo() != "いい会社" {
		t.Errorf("Memo() = %q, want %q", company.Memo(), "いい会社")
	}

	company.UpdateMemo("")
	if company.Memo() != "" {
		t.Errorf("Memo() should be empty after clearing, got %q", company.Memo())
	}
}
