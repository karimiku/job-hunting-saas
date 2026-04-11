package entity

import (
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestCompanyName(t *testing.T, raw string) value.CompanyName {
	t.Helper()
	name, err := value.NewCompanyName(raw)
	if err != nil {
		t.Fatalf("NewCompanyName failed: %v", err)
	}
	return name
}

func TestNewCompany(t *testing.T) {
	userID := NewUserID()
	name := newTestCompanyName(t, "トヨタ自動車")

	t.Run("valid company", func(t *testing.T) {
		company := NewCompany(userID, name)
		if company.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if company.UserID() != userID {
			t.Errorf("UserID() = %v, want %v", company.UserID(), userID)
		}
		if company.Name().String() != "トヨタ自動車" {
			t.Errorf("Name() = %q, want %q", company.Name().String(), "トヨタ自動車")
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
}

func TestCompany_Rename(t *testing.T) {
	userID := NewUserID()
	name := newTestCompanyName(t, "旧社名")
	company := NewCompany(userID, name)

	newName := newTestCompanyName(t, "新社名")
	company.Rename(newName)

	if company.Name().String() != "新社名" {
		t.Errorf("Name() = %q, want %q", company.Name().String(), "新社名")
	}
}

func TestCompany_UpdateMemo(t *testing.T) {
	userID := NewUserID()
	name := newTestCompanyName(t, "テスト株式会社")
	company := NewCompany(userID, name)

	company.UpdateMemo("いい会社")
	if company.Memo() != "いい会社" {
		t.Errorf("Memo() = %q, want %q", company.Memo(), "いい会社")
	}

	company.UpdateMemo("")
	if company.Memo() != "" {
		t.Errorf("Memo() should be empty after clearing, got %q", company.Memo())
	}
}
