package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewCompanyAlias(t *testing.T) {
	userID := NewID()
	companyID := NewID()

	t.Run("valid alias", func(t *testing.T) {
		alias, err := NewCompanyAlias(userID, companyID, "ソニー")
		if err != nil {
			t.Fatalf("NewCompanyAlias should succeed, but got error: %v", err)
		}
		if alias.ID() == uuid.Nil {
			t.Error("ID should not be nil")
		}
		if alias.UserID() != userID {
			t.Errorf("UserID() = %v, want %v", alias.UserID(), userID)
		}
		if alias.CompanyID() != companyID {
			t.Errorf("CompanyID() = %v, want %v", alias.CompanyID(), companyID)
		}
		if alias.Alias() != "ソニー" {
			t.Errorf("Alias() = %q, want %q", alias.Alias(), "ソニー")
		}
	})

	t.Run("empty alias", func(t *testing.T) {
		_, err := NewCompanyAlias(userID, companyID, "")
		if err == nil {
			t.Error("NewCompanyAlias with empty alias should return error")
		}
	})

	t.Run("whitespace alias", func(t *testing.T) {
		_, err := NewCompanyAlias(userID, companyID, "   ")
		if err == nil {
			t.Error("NewCompanyAlias with whitespace alias should return error")
		}
	})
}
