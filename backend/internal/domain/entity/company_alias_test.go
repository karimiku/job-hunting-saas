package entity

import (
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestAlias(t *testing.T, raw string) value.Alias {
	t.Helper()
	alias, err := value.NewAlias(raw)
	if err != nil {
		t.Fatalf("NewAlias failed: %v", err)
	}
	return alias
}

func TestNewCompanyAlias(t *testing.T) {
	userID := NewUserID()
	companyID := NewCompanyID()
	alias := newTestAlias(t, "ソニー")

	t.Run("valid alias", func(t *testing.T) {
		ca := NewCompanyAlias(userID, companyID, alias)
		if ca.ID().IsZero() {
			t.Error("ID should not be zero")
		}
		if ca.UserID() != userID {
			t.Errorf("UserID() = %v, want %v", ca.UserID(), userID)
		}
		if ca.CompanyID() != companyID {
			t.Errorf("CompanyID() = %v, want %v", ca.CompanyID(), companyID)
		}
		if ca.Alias().String() != "ソニー" {
			t.Errorf("Alias() = %q, want %q", ca.Alias().String(), "ソニー")
		}
	})
}
