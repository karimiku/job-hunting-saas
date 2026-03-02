package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

func TestUpdateCompany_Success(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	saveCalled := false
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Company) error {
			saveCalled = true
			return nil
		},
	}

	uc := NewUpdateCompany(repo)
	out, err := uc.Execute(context.Background(), UpdateCompanyInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      "新しい社名",
		Memo:      "新しいメモ",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Company.Name().String() != "新しい社名" {
		t.Errorf("Name = %q, want %q", out.Company.Name().String(), "新しい社名")
	}
	if out.Company.Memo() != "新しいメモ" {
		t.Errorf("Memo = %q, want %q", out.Company.Memo(), "新しいメモ")
	}
	if !saveCalled {
		t.Error("Save should be called")
	}
}

func TestUpdateCompany_EmptyName(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewUpdateCompany(repo)
	_, err := uc.Execute(context.Background(), UpdateCompanyInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Name:      "",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, value.ErrCompanyNameEmpty) {
		t.Errorf("error = %v, want ErrCompanyNameEmpty", err)
	}
}

func TestUpdateCompany_NotFound(t *testing.T) {
	repo := &mockCompanyRepo{}

	uc := NewUpdateCompany(repo)
	_, err := uc.Execute(context.Background(), UpdateCompanyInput{
		UserID:    entity.NewUserID(),
		CompanyID: entity.NewCompanyID(),
		Name:      "株式会社テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("error = %v, want ErrNotFound", err)
	}
}

func TestUpdateCompany_SaveError(t *testing.T) {
	userID := entity.NewUserID()
	existing := newTestCompany(t, userID)
	saveErr := errors.New("db write failed")
	repo := &mockCompanyRepo{
		findByIDFn: func(_ context.Context, _ entity.UserID, _ entity.CompanyID) (*entity.Company, error) {
			return existing, nil
		},
		saveFn: func(_ context.Context, _ *entity.Company) error {
			return saveErr
		},
	}

	uc := NewUpdateCompany(repo)
	_, err := uc.Execute(context.Background(), UpdateCompanyInput{
		UserID:    userID,
		CompanyID: existing.ID(),
		Name:      "株式会社テスト",
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, saveErr) {
		t.Errorf("error = %v, want saveErr", err)
	}
}
