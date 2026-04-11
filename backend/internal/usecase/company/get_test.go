package company

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func newTestCompany(t *testing.T, userID entity.UserID) *entity.Company {
	t.Helper()
	name, err := value.NewCompanyName("株式会社テスト")
	if err != nil {
		t.Fatalf("NewCompanyName failed: %v", err)
	}
	return entity.NewCompany(userID, name)
}

func TestGet_Found(t *testing.T) {
	userID := entity.NewUserID()
	expected := newTestCompany(t, userID)
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return expected, nil
		},
	}

	uc := NewGet(repo)
	out, err := uc.Execute(context.Background(), GetInput{
		UserID:    userID,
		CompanyID: expected.ID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company != expected {
		t.Error("Company should be the expected company")
	}
}

func TestGet_NotFound(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestGet_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return nil, dbErr
		},
	}

	uc := NewGet(repo)
	_, err := uc.Execute(context.Background(), GetInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
