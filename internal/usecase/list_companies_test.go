package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

func TestListCompanies_Multiple(t *testing.T) {
	userID := entity.NewUserID()
	c1 := newTestCompany(t, userID)
	c2 := newTestCompany(t, userID)
	repo := &mockCompanyRepo{
		listByUserFn: func(_ context.Context, _ entity.UserID) ([]*entity.Company, error) {
			return []*entity.Company{c1, c2}, nil
		},
	}

	uc := NewListCompanies(repo)
	out, err := uc.Execute(context.Background(), ListCompaniesInput{
		UserID: userID,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Companies) != 2 {
		t.Fatalf("len(Companies) = %d, want 2", len(out.Companies))
	}
}

func TestListCompanies_Empty(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewListCompanies(repo)
	out, err := uc.Execute(context.Background(), ListCompaniesInput{
		UserID: entity.NewUserID(),
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Companies) != 0 {
		t.Fatalf("len(Companies) = %d, want 0", len(out.Companies))
	}
}

func TestListCompanies_DBError(t *testing.T) {
	dbErr := errors.New("db connection failed")
	repo := &mockCompanyRepo{
		listByUserFn: func(_ context.Context, _ entity.UserID) ([]*entity.Company, error) {
			return nil, dbErr
		},
	}

	uc := NewListCompanies(repo)
	_, err := uc.Execute(context.Background(), ListCompaniesInput{
		UserID: entity.NewUserID(),
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dbErr) {
		t.Errorf("error = %v, want dbErr", err)
	}
}
